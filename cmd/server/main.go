package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	_ "modernc.org/sqlite"

	"github.com/lingvo-ai/lingvo/internal/api"
	"github.com/lingvo-ai/lingvo/internal/bot"
	"github.com/lingvo-ai/lingvo/internal/db"
)

func main() {
	// Logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	// Config from env
	botToken := os.Getenv("BOT_TOKEN")
	geminiKey := os.Getenv("GEMINI_API_KEY")
	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		dbPath = "lingvo.db"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Init database
	database, err := sqlx.Connect("sqlite", dbPath)
	if err != nil {
		sugar.Fatalw("db connect", "error", err)
	}
	defer database.Close()

	db.Migrate(database, sugar)

	// Init router
	r := gin.Default()

	// API routes (includes /api/health)
	api.RegisterRoutes(r, database, geminiKey, botToken, sugar)

	// No-cache middleware for frontend assets
	r.Use(func(c *gin.Context) {
		if !c.IsAborted() {
			c.Header("Cache-Control", "no-store")
		}
		c.Next()
	})

	// Serve static frontend files
	r.Static("/assets", "./web/dist/assets")
	r.StaticFile("/favicon.svg", "./web/dist/favicon.svg")
	r.StaticFile("/icons.svg", "./web/dist/icons.svg")

	// SPA fallback — serve index.html for all unmatched routes
	r.NoRoute(func(c *gin.Context) {
		c.File("./web/dist/index.html")
	})

	// Start bot (long-polling, blocking in goroutine)
	if botToken != "" {
		go bot.Start(database, botToken, sugar)
	}

	sugar.Infow("starting server", "port", port)
	log.Fatal(r.Run(":" + port))
}

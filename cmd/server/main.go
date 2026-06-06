package main

import (
	"bufio"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	_ "modernc.org/sqlite"

	"github.com/lingvo-ai/lingvo/internal/api"
	"github.com/lingvo-ai/lingvo/internal/bot"
	"github.com/lingvo-ai/lingvo/internal/db"
)

func loadEnv(path string) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		if os.Getenv(key) == "" {
			os.Setenv(key, val)
		}
	}
}

func main() {
	loadEnv(".env")

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

	// API routes
	api.RegisterRoutes(r, database, geminiKey, botToken, sugar)

	// Serve frontend (SPA — index.html for all non-API routes)
	r.Static("/assets", "./web/dist/assets")
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

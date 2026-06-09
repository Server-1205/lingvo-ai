package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	_ "modernc.org/sqlite"

	"github.com/lingvo-ai/lingvo/internal/ai"
	"github.com/lingvo-ai/lingvo/internal/api"
	"github.com/lingvo-ai/lingvo/internal/bot"
	"github.com/lingvo-ai/lingvo/internal/db"
	"github.com/lingvo-ai/lingvo/internal/tts"
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
	openAIKey := os.Getenv("OPENAI_API_KEY")
	openAIBaseURL := os.Getenv("OPENAI_BASE_URL")
	openAIModel := os.Getenv("OPENAI_MODEL")
	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		dbPath = "lingvo.db"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	aiQueueEnabled := os.Getenv("AI_QUEUE_ENABLED") == "true"
	devMode := os.Getenv("DEV_MODE") == "true"
	adminIDs := os.Getenv("ADMIN_IDS")
	webappURL := os.Getenv("WEBAPP_URL")
	if webappURL == "" {
		webappURL = "https://t.me/lingvo_ai_bot/app"
	}

	// Init database
	database, err := sqlx.Connect("sqlite", dbPath)
	if err != nil {
		sugar.Fatalw("db connect", "error", err)
	}
	defer database.Close()

	db.Migrate(database, sugar)

	// Init AI client
	var aiClient *ai.Client
	if geminiKey != "" {
		var err error
		aiClient, err = ai.NewClient(context.Background(), geminiKey, openAIKey, openAIBaseURL, openAIModel, sugar)
		if err != nil {
			sugar.Warnw("failed to init AI client, AI features disabled", "error", err)
		} else {
			sugar.Info("AI client initialized")
			if aiQueueEnabled {
				aiClient.EnableQueue(context.Background(), sugar)
				sugar.Info("AI priority queue enabled")
			}
		}
	} else {
		sugar.Warn("GEMINI_API_KEY not set, AI features disabled")
	}

	// Init router
	r := gin.Default()

	ttsVoiceUz := os.Getenv("TTS_VOICE_UZ")
	if ttsVoiceUz == "" {
		ttsVoiceUz = tts.DefaultVoices.Uz
	}
	ttsVoiceRu := os.Getenv("TTS_VOICE_RU")
	if ttsVoiceRu == "" {
		ttsVoiceRu = tts.DefaultVoices.Ru
	}
	ttsClient := tts.NewClient(ttsVoiceUz, ttsVoiceRu, sugar)
	if !ttsClient.IsAvailable() {
		sugar.Warn("[tts] edge-tts not installed — TTS endpoint will return 500")
	}

	api.RegisterRoutes(r, database, botToken, sugar, adminIDs, devMode, aiClient, ttsClient)

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
		var parsedAdminIDs []int64
		if adminIDs != "" {
			parts := strings.Split(adminIDs, ",")
			for _, p := range parts {
				p = strings.TrimSpace(p)
				if p == "" {
					continue
				}
				var id int64
				if _, err := fmt.Sscanf(p, "%d", &id); err == nil {
					parsedAdminIDs = append(parsedAdminIDs, id)
				}
			}
		}
		go bot.Start(database, botToken, webappURL, sugar, parsedAdminIDs, aiClient)
	}

	sugar.Infow("starting server", "port", port)
	log.Fatal(r.Run(":" + port))
}

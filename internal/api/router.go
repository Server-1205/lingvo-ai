package api

import (
	"context"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/lingvo-ai/lingvo/internal/ai"
	"github.com/lingvo-ai/lingvo/internal/middleware"
)

func parseAdminIDs(raw string) []int64 {
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	ids := make([]int64, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		var id int64
		if _, err := fmt.Sscanf(p, "%d", &id); err == nil {
			ids = append(ids, id)
		}
	}
	return ids
}

func RegisterRoutes(r *gin.Engine, db *sqlx.DB, geminiKey, openAIKey, openAIBaseURL, openAIModel, botToken string, sugar *zap.SugaredLogger, aiQueueEnabled bool, adminIDList string, devMode bool) {
	api := r.Group("/api")
	api.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	authMw := middleware.AuthMiddleware(botToken, db)
	adminIDs := parseAdminIDs(adminIDList)
	rateMw := middleware.RateLimitMiddleware(db, adminIDs, devMode)

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

	protected := api.Group("")
	protected.Use(authMw)
	{
		protected.POST("/chat", rateMw, chatHandler(db, aiClient, sugar))
		protected.POST("/chat/stream", rateMw, chatStreamHandler(db, aiClient, sugar))
		protected.POST("/grammar", rateMw, grammarHandler(db, aiClient, sugar))
		protected.GET("/vocab", rateMw, vocabListHandler(db))
		protected.POST("/vocab", vocabAddHandler(db, aiClient, sugar))
		protected.DELETE("/vocab/:id", vocabDeleteHandler(db, sugar))
		protected.POST("/vocab/lookup", vocabLookupHandler(aiClient, sugar))
		protected.GET("/vocab/export", vocabExportHandler(db, sugar))
		protected.GET("/vocab/review", vocabReviewHandler(db))
		protected.POST("/vocab/review", vocabReviewSubmitHandler(db, sugar))
		protected.POST("/quiz", rateMw, quizHandler(db, aiClient, sugar))
		protected.POST("/level", levelTestHandler(db, aiClient, sugar))
		protected.POST("/level/save", levelSaveHandler(db))
		protected.GET("/progress", progressHandler(db))
		protected.GET("/progress/history", progressHistoryHandler(db, sugar))
		protected.GET("/subscription", subscriptionHandler(db))
		protected.POST("/create-invoice", invoiceHandler(botToken, sugar))
		protected.GET("/errors/history", rateMw, errorsHistoryHandler(db, sugar))
		protected.GET("/errors/stats", rateMw, errorsStatsHandler(db, sugar))
		protected.POST("/ielts/writing", rateMw, ieltsWritingHandler(db, aiClient, sugar))
		protected.GET("/ielts/speaking", rateMw, ieltsSpeakingGenerateHandler(db, aiClient, sugar))
		protected.POST("/ielts/speaking", rateMw, ieltsSpeakingEvaluateHandler(db, aiClient, sugar))
		protected.GET("/ielts/reading", rateMw, ieltsReadingGenerateHandler(db, aiClient, sugar))
		protected.POST("/ielts/reading/submit", rateMw, ieltsReadingEvaluateHandler(db, aiClient, sugar))
		protected.GET("/ielts/scores", rateMw, ieltsScoresHandler(db, sugar))
	}
}

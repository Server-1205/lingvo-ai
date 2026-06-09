package api

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/lingvo-ai/lingvo/internal/ai"
	"github.com/lingvo-ai/lingvo/internal/cache"
	"github.com/lingvo-ai/lingvo/internal/middleware"
	"github.com/lingvo-ai/lingvo/internal/tts"
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

func RegisterRoutes(r *gin.Engine, db *sqlx.DB, botToken string, sugar *zap.SugaredLogger, adminIDList string, devMode bool, aiClient *ai.Client, ttsClient tts.Synthesizer) {
	api := r.Group("/api")
	api.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	authMw := middleware.AuthMiddleware(botToken, db, devMode)
	adminIDs := parseAdminIDs(adminIDList)
	rateMw := middleware.RateLimitMiddleware(db, adminIDs, devMode)
	apiCache := cache.New(30 * time.Second)

	if aiClient != nil {
		go func() {
			normalizeVocabCase(context.Background(), db, sugar)
			translateMissingVocab(context.Background(), db, aiClient, sugar)
		}()
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
		protected.POST("/daily", rateMw, dailyHandler(db, aiClient, sugar))
		protected.GET("/progress", progressHandler(db, apiCache))
		protected.GET("/progress/history", progressHistoryHandler(db, sugar, apiCache))
		protected.GET("/subscription", subscriptionHandler(db, apiCache))
		protected.POST("/create-invoice", invoiceHandler(botToken, sugar))
		protected.GET("/errors/history", rateMw, errorsHistoryHandler(db, sugar))
		protected.GET("/errors/stats", rateMw, errorsStatsHandler(db, sugar, apiCache))
		protected.POST("/ielts/writing", rateMw, ieltsWritingHandler(db, aiClient, sugar))
		protected.GET("/ielts/speaking", rateMw, ieltsSpeakingGenerateHandler(db, aiClient, sugar))
		protected.POST("/ielts/speaking", rateMw, ieltsSpeakingEvaluateHandler(db, aiClient, sugar))
		protected.GET("/ielts/reading", rateMw, ieltsReadingGenerateHandler(db, aiClient, sugar))
		protected.POST("/ielts/reading/submit", rateMw, ieltsReadingEvaluateHandler(db, aiClient, sugar))
		protected.GET("/ielts/scores", rateMw, ieltsScoresHandler(db, sugar))
		protected.GET("/tts", ttsHandler(ttsClient, sugar))
	}
}

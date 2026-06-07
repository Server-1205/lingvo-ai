package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/lingvo-ai/lingvo/internal/ai"
	"github.com/lingvo-ai/lingvo/internal/db"
	"github.com/lingvo-ai/lingvo/internal/models"
)

func quizHandler(database *sqlx.DB, aiClient *ai.Client, sugar *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")

		if aiClient == nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "ai_service_unavailable"})
			return
		}

		var req models.QuizRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request"})
			return
		}

		if req.Count == 0 || req.Count > 10 {
			req.Count = 5
		}
		if req.Topic == "" {
			req.Topic = "general"
		}

		lang, _ := c.Get("lang")
		lng, _ := lang.(string)

		prompt := ai.BuildQuizPrompt(req.Topic, req.Count, lng)

		raw, err := aiClient.GenerateLite(c.Request.Context(), prompt)
		if err != nil {
			sugar.Errorw("ai quiz error", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ai_service_unavailable"})
			return
		}

		raw = cleanJSON(raw)

		var resp models.QuizResponse
		if err := json.Unmarshal([]byte(raw), &resp); err != nil {
			sugar.Errorw("parse quiz response", "error", err, "raw", raw)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "parse_error"})
			return
		}

		userID, _ := c.Get("user_id")
		uid, _ := userID.(int)
		today := time.Now().UTC().Format("2006-01-02")
		_ = db.IncrementProgress(c.Request.Context(), database, uid, today, "quizzes_taken")

		c.JSON(http.StatusOK, resp)
	}
}

package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/lingvo-ai/lingvo/internal/ai"
	"github.com/lingvo-ai/lingvo/internal/db"
	"github.com/lingvo-ai/lingvo/internal/models"
)

func chatHandler(database *sqlx.DB, aiClient *ai.Client, sugar *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		if aiClient == nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "ai_service_unavailable"})
			return
		}

		var req models.ChatRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request"})
			return
		}

		userID, _ := c.Get("user_id")
		uid, _ := userID.(int)
		lang, _ := c.Get("lang")
		level, _ := c.Get("level")
		isPremium, _ := c.Get("is_premium")
		dailyUsed, _ := c.Get("daily_used")
		dailyLimit, _ := c.Get("daily_limit")

		lvl, _ := level.(string)
		lng, _ := lang.(string)
		used, _ := dailyUsed.(int)
		limit, _ := dailyLimit.(int)
		premium, _ := isPremium.(bool)

		prompt := ai.BuildChatPrompt(lvl, lng, req.Text)

		aiResp, err := aiClient.Chat(c.Request.Context(), prompt)
		if err != nil {
			sugar.Errorw("ai chat error", "error", err, "telegram_id", c.GetInt64("telegram_id"))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ai_service_unavailable"})
			return
		}

		today := time.Now().UTC().Format("2006-01-02")
		if err := db.IncrementMessageCount(c.Request.Context(), database, uid, today); err != nil {
			sugar.Errorw("increment message count", "error", err)
		}

		var corrections []models.Correction
		if aiResp != nil && aiResp.Corrections != nil {
			for _, cr := range aiResp.Corrections {
				corrections = append(corrections, models.Correction{
					Original:      cr.Original,
					Corrected:     cr.Corrected,
					ExplanationUz: cr.ExplanationUz,
					ExplanationRu: cr.ExplanationRu,
					Type:          cr.Type,
				})
			}
		}

		reply := ""
		if aiResp != nil {
			reply = aiResp.Reply
		}

		c.JSON(http.StatusOK, models.ChatResponse{
			Reply:       reply,
			Corrections: corrections,
			Usage: models.Usage{
				DailyUsed:  used + 1,
				DailyLimit: limit,
				IsPremium:  premium,
			},
		})
	}
}

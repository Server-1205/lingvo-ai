package api

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/lingvo-ai/lingvo/internal/ai"
)

func dailyHandler(database *sqlx.DB, aiClient *ai.Client, sugar *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")

		if aiClient == nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "ai_service_unavailable"})
			return
		}

		userID, _ := c.Get("user_id")
		uid, _ := userID.(int)
		lang, _ := c.Get("lang")
		level, _ := c.Get("level")
		lvl, _ := level.(string)
		lng, _ := lang.(string)

		sugar.Infow("[daily] generating lesson", "user_id", uid, "level", lvl, "lang", lng)

		prompt := ai.BuildDailyLessonPrompt(lvl, lng)
		raw, err := aiClient.Generate(c.Request.Context(), prompt)
		if err != nil {
			sugar.Errorw("[daily] AI generate failed", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ai_service_unavailable"})
			return
		}

		raw = cleanJSON(raw)
		var lesson ai.DailyLessonResponse
		if err := json.Unmarshal([]byte(raw), &lesson); err != nil {
			sugar.Errorw("[daily] parse lesson failed", "error", err, "raw", raw)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "parse_error"})
			return
		}

		if lesson.Topic == "" {
			sugar.Errorw("[daily] empty topic in response", "raw", raw)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "parse_error"})
			return
		}

		sugar.Infow("[daily] lesson generated", "user_id", uid, "topic", lesson.Topic)

		c.JSON(http.StatusOK, lesson)
	}
}

package api

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/lingvo-ai/lingvo/internal/ai"
	"github.com/lingvo-ai/lingvo/internal/models"
)

func levelTestHandler(database *sqlx.DB, aiClient *ai.Client, sugar *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.LevelRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request"})
			return
		}

		lang, _ := c.Get("lang")
		lng, _ := lang.(string)

		prompt := ai.BuildLevelTestPrompt(lng)

		raw, err := aiClient.Generate(c.Request.Context(), prompt)
		if err != nil {
			sugar.Errorw("ai level test error", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ai_service_unavailable"})
			return
		}

		raw = cleanJSON(raw)

		var levelResp models.LevelResponse
		if err := json.Unmarshal([]byte(raw), &levelResp); err != nil {
			sugar.Errorw("parse level response", "error", err, "raw", raw)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "parse_error"})
			return
		}

		c.JSON(http.StatusOK, levelResp)
	}
}

func levelSaveHandler(database *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Level string `json:"level" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request"})
			return
		}

		valid := map[string]bool{"a1": true, "a2": true, "b1": true, "b2": true, "c1": true}
		if !valid[req.Level] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_level"})
			return
		}

		userID, _ := c.Get("user_id")
		uid, _ := userID.(int)

		_, err := database.ExecContext(c.Request.Context(),
			"UPDATE users SET level = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?",
			req.Level, uid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "ok", "level": req.Level})
	}
}

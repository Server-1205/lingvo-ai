package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/lingvo-ai/lingvo/internal/ai"
	"github.com/lingvo-ai/lingvo/internal/db"
	"github.com/lingvo-ai/lingvo/internal/models"
)

type sseEvent struct {
	Type string      `json:"type"`
	Data interface{} `json:"data,omitempty"`
}

func chatStreamHandler(database *sqlx.DB, aiClient *ai.Client, sugar *zap.SugaredLogger) gin.HandlerFunc {
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

		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")

		sugar.Infow("streaming session started", "user_id", uid)

		streamCh := aiClient.ChatStream(c.Request.Context(), prompt)

		var finalResult *models.AIResponse

	loop:
		for {
			select {
			case evt, ok := <-streamCh:
				if !ok {
					break loop
				}
				switch evt.Type {
				case ai.EventToken:
					var token string
					if err := json.Unmarshal(evt.Data, &token); err != nil {
						sugar.Errorw("failed to unmarshal token", "error", err)
						continue
					}
					sugar.Debugw("streaming token", "user_id", uid, "len", len(token))
					writeSSE(c, sseEvent{Type: "token", Data: token})
				case ai.EventResult:
					var result models.AIResponse
					if err := json.Unmarshal(evt.Data, &result); err != nil {
						sugar.Errorw("failed to unmarshal result", "error", err)
						continue
					}
					finalResult = &result
				}
			case <-c.Request.Context().Done():
				sugar.Warnw("client disconnected", "user_id", uid)
				return
			}
		}

		today := time.Now().UTC().Format("2006-01-02")
		if err := db.IncrementMessageCount(c.Request.Context(), database, uid, today); err != nil {
			sugar.Errorw("increment message count", "error", err)
		}

		var corrections []models.Correction
		if finalResult != nil && finalResult.Corrections != nil {
			for _, cr := range finalResult.Corrections {
				corrections = append(corrections, models.Correction{
					Original:      cr.Original,
					Corrected:     cr.Corrected,
					ExplanationUz: cr.ExplanationUz,
					ExplanationRu: cr.ExplanationRu,
					Type:          cr.Type,
				})
			}
		}

		writeSSE(c, sseEvent{Type: "corrections", Data: corrections})
		writeSSE(c, sseEvent{
			Type: "usage",
			Data: models.Usage{
				DailyUsed:  used + 1,
				DailyLimit: limit,
				IsPremium:  premium,
			},
		})
		writeSSE(c, sseEvent{Type: "done"})

		sugar.Infow("streaming session completed", "user_id", uid)
	}
}

func writeSSE(c *gin.Context, evt sseEvent) {
	data, err := json.Marshal(evt)
	if err != nil {
		return
	}
	_, err = fmt.Fprintf(c.Writer, "data: %s\n\n", data)
	if err != nil {
		return
	}
	c.Writer.Flush()
}

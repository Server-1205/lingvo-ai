package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"

	"github.com/lingvo-ai/lingvo/internal/db"
)

func RateLimitMiddleware(database *sqlx.DB, adminIDs []int64, devMode bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		uid, ok := userID.(int)
		if !ok {
			c.AbortWithStatusJSON(500, gin.H{"error": "internal_error"})
			return
		}

		if devMode {
			c.Set("is_premium", true)
			c.Set("daily_used", 0)
			c.Set("daily_limit", 0)
			c.Next()
			return
		}

		tgID, _ := c.Get("telegram_id")
		if tgIDInt, ok := tgID.(int64); ok {
			for _, adminID := range adminIDs {
				if adminID == tgIDInt {
					c.Set("is_premium", true)
					c.Set("daily_used", 0)
					c.Set("daily_limit", 0)
					c.Next()
					return
				}
			}
		}

		sub, err := db.GetActiveSubscription(c.Request.Context(), database, uid)
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": "internal_error"})
			return
		}

		if sub != nil {
			c.Set("is_premium", true)
			c.Set("daily_used", 0)
			c.Set("daily_limit", 0)
			c.Next()
			return
		}

		today := time.Now().UTC().Format("2006-01-02")
		count, err := db.GetMessageCount(c.Request.Context(), database, uid, today)
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": "internal_error"})
			return
		}

		limit := 10
		if count >= limit {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":        "daily_limit_exceeded",
				"message_uz":   "Бугунги лимит тугади. Чексизга обуна бўлинг!",
				"message_ru":   "Дневной лимит исчерпан. Оформите подписку!",
				"daily_used":   count,
				"daily_limit":  limit,
				"is_premium":   false,
				"premium_link": "https://t.me/lingvo_ai_bot/app",
			})
			return
		}

		c.Set("is_premium", false)
		c.Set("daily_used", count)
		c.Set("daily_limit", limit)
		c.Next()
	}
}

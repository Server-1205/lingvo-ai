package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"

	"github.com/lingvo-ai/lingvo/internal/cache"
	"github.com/lingvo-ai/lingvo/internal/db"
	"github.com/lingvo-ai/lingvo/internal/models"
)

func subscriptionHandler(database *sqlx.DB, cch *cache.Cache) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		uid, _ := userID.(int)

		cacheKey := "subscription:" + strconv.Itoa(uid)
		if data, ok := cch.Get(cacheKey); ok {
			c.Data(http.StatusOK, "application/json; charset=utf-8", data.([]byte))
			c.Abort()
			return
		}

		sub, err := db.GetActiveSubscription(c.Request.Context(), database, uid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error"})
			return
		}

		var resp models.SubscriptionResponse
		if sub == nil {
			resp = models.SubscriptionResponse{Active: false}
		} else {
			resp = models.SubscriptionResponse{
				Active:    true,
				Plan:      sub.Plan,
				ExpiresAt: sub.ExpiresAt.Format("2006-01-02 15:04:05"),
			}
		}

		if raw, err := json.Marshal(resp); err == nil {
			cch.Set(cacheKey, raw, 60*time.Second)
		}
		c.JSON(http.StatusOK, resp)
	}
}

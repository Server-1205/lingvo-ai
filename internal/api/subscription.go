package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"

	"github.com/lingvo-ai/lingvo/internal/db"
	"github.com/lingvo-ai/lingvo/internal/models"
)

func subscriptionHandler(database *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		uid, _ := userID.(int)

		sub, err := db.GetActiveSubscription(c.Request.Context(), database, uid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error"})
			return
		}

		if sub == nil {
			c.JSON(http.StatusOK, models.SubscriptionResponse{Active: false})
			return
		}

		c.JSON(http.StatusOK, models.SubscriptionResponse{
			Active:    true,
			Plan:      sub.Plan,
			ExpiresAt: sub.ExpiresAt.Format("2006-01-02 15:04:05"),
		})
	}
}

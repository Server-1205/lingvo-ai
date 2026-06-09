package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/lingvo-ai/lingvo/internal/db"
	"github.com/lingvo-ai/lingvo/internal/models"
)

func errorsHistoryHandler(database *sqlx.DB, sugar *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")

		isPremium, _ := c.Get("is_premium")
		premium, _ := isPremium.(bool)
		if !premium {
			c.JSON(http.StatusForbidden, gin.H{"error": "premium_required"})
			return
		}

		userID, _ := c.Get("user_id")
		uid, _ := userID.(int)

		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		if page < 1 {
			page = 1
		}
		perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
		if perPage < 1 {
			perPage = 20
		}
		if perPage > 100 {
			perPage = 100
		}
		category := c.DefaultQuery("category", "")

		entries, total, err := db.GetErrorHistory(c.Request.Context(), database, uid, perPage, (page-1)*perPage, category)
		if err != nil {
			sugar.Errorw("error history query failed", "user_id", uid, "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error"})
			return
		}

		sugar.Debugw("error history requested", "user_id", uid, "total", total, "page", page)

		c.JSON(http.StatusOK, models.ErrorHistoryResponse{
			Entries: entries,
			Total:   total,
		})
	}
}

func errorsStatsHandler(database *sqlx.DB, sugar *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")

		isPremium, _ := c.Get("is_premium")
		premium, _ := isPremium.(bool)
		if !premium {
			c.JSON(http.StatusForbidden, gin.H{"error": "premium_required"})
			return
		}

		userID, _ := c.Get("user_id")
		uid, _ := userID.(int)

		days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))
		if days < 1 {
			days = 30
		}
		if days > 365 {
			days = 365
		}

		stats, err := db.GetErrorStats(c.Request.Context(), database, uid, days)
		if err != nil {
			sugar.Errorw("error stats query failed", "user_id", uid, "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error"})
			return
		}

		sugar.Debugw("error stats requested", "user_id", uid, "total_errors", stats.TotalErrors)

		c.JSON(http.StatusOK, stats)
	}
}

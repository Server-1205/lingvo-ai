package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/lingvo-ai/lingvo/internal/cache"
	"github.com/lingvo-ai/lingvo/internal/db"
	"github.com/lingvo-ai/lingvo/internal/models"
)

func progressHistoryHandler(database *sqlx.DB, sugar *zap.SugaredLogger, cch *cache.Cache) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")
		userID, _ := c.Get("user_id")
		uid, _ := userID.(int)

		days := 7
		if d := c.Query("days"); d != "" {
			if parsed, err := strconv.Atoi(d); err == nil && parsed > 0 && parsed <= 30 {
				days = parsed
			}
		}

		cacheKey := "progress_history:" + strconv.Itoa(uid) + ":" + strconv.Itoa(days)
		if data, ok := cch.Get(cacheKey); ok {
			c.Data(http.StatusOK, "application/json; charset=utf-8", data.([]byte))
			c.Abort()
			return
		}

		sugar.Debugw("progress history requested", "user_id", uid, "days", days)

		entries, err := db.GetProgressHistory(c.Request.Context(), database, uid, days)
		if err != nil {
			sugar.Errorw("progress history db error", "user_id", uid, "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error"})
			return
		}
		if entries == nil {
			entries = []models.DailyProgressEntry{}
		}

		sugar.Debugw("progress history response", "user_id", uid, "entries", len(entries))
		resp := models.ProgressHistoryResponse{Entries: entries}
		if raw, err := json.Marshal(resp); err == nil {
			cch.SetDefault(cacheKey, raw)
		}
		c.JSON(http.StatusOK, resp)
	}
}

func progressHandler(database *sqlx.DB, cch *cache.Cache) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		uid, _ := userID.(int)
		level, _ := c.Get("level")
		lvl, _ := level.(string)

		cacheKey := "progress:" + strconv.Itoa(uid)
		if data, ok := cch.Get(cacheKey); ok {
			c.Data(http.StatusOK, "application/json; charset=utf-8", data.([]byte))
			c.Abort()
			return
		}

		today := time.Now().UTC().Format("2006-01-02")

		progress, err := db.GetDailyProgress(c.Request.Context(), database, uid, today)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error"})
			return
		}

		streak, err := db.GetStreakDays(c.Request.Context(), database, uid)
		if err != nil {
			streak = 0
		}

		resp := models.ProgressResponse{
			MessagesSent: progress["messages_sent"],
			WordsLearned: progress["words_learned"],
			QuizzesTaken: progress["quizzes_taken"],
			StreakDays:   streak,
			Level:        lvl,
		}
		if raw, err := json.Marshal(resp); err == nil {
			cch.SetDefault(cacheKey, raw)
		}
		c.JSON(http.StatusOK, resp)
	}
}

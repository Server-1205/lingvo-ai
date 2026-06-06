package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/lingvo-ai/lingvo/internal/db"
	"github.com/lingvo-ai/lingvo/internal/models"
)

func progressHistoryHandler(database *sqlx.DB, sugar *zap.SugaredLogger) gin.HandlerFunc {
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

		sugar.Debugw("progress history requested", "user_id", uid, "days", days)

		entries, err := db.GetProgressHistory(c.Request.Context(), database, uid, days)
		if err != nil {
			sugar.Errorw("progress history db error", "user_id", uid, "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error"})
			return
		}

		sugar.Debugw("progress history response", "user_id", uid, "entries", len(entries))
		c.JSON(http.StatusOK, models.ProgressHistoryResponse{Entries: entries})
	}
}

func progressHandler(database *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		uid, _ := userID.(int)
		level, _ := c.Get("level")
		lvl, _ := level.(string)

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

		c.JSON(http.StatusOK, models.ProgressResponse{
			MessagesSent: progress["messages_sent"],
			WordsLearned: progress["words_learned"],
			QuizzesTaken: progress["quizzes_taken"],
			StreakDays:   streak,
			Level:        lvl,
		})
	}
}

package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"

	"github.com/lingvo-ai/lingvo/internal/db"
	"github.com/lingvo-ai/lingvo/internal/models"
)

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

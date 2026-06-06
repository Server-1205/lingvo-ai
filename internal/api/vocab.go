package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/lingvo-ai/lingvo/internal/ai"
	"github.com/lingvo-ai/lingvo/internal/db"
	"github.com/lingvo-ai/lingvo/internal/models"
)

type reviewRequest struct {
	Quality int `json:"quality" binding:"min=0,max=5"`
}

type reviewResponse struct {
	NextReview string `json:"next_review"`
}

func vocabListHandler(database *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		uid, _ := userID.(int)

		words, err := db.GetVocabulary(c.Request.Context(), database, uid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error"})
			return
		}

		if words == nil {
			words = []models.VocabWord{}
		}

		c.JSON(http.StatusOK, words)
	}
}

func vocabAddHandler(database *sqlx.DB, sugar *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.AddVocabRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request"})
			return
		}

		userID, _ := c.Get("user_id")
		uid, _ := userID.(int)

		if err := db.AddVocabulary(c.Request.Context(), database, uid, req.Word, req.Translation, req.Example, req.Level); err != nil {
			sugar.Errorw("add vocab", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"status": "ok"})
	}
}

func vocabDeleteHandler(database *sqlx.DB, sugar *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_id"})
			return
		}

		userID, _ := c.Get("user_id")
		uid, _ := userID.(int)

		if err := db.DeleteVocabulary(c.Request.Context(), database, uid, id); err != nil {
			sugar.Errorw("delete vocab", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}

func vocabReviewHandler(database *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		uid, _ := userID.(int)

		limitStr := c.DefaultQuery("limit", "20")
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit < 1 {
			limit = 20
		}

		words, err := db.GetDueWords(c.Request.Context(), database, uid, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error"})
			return
		}

		if words == nil {
			words = []models.VocabWord{}
		}

		c.JSON(http.StatusOK, words)
	}
}

func vocabReviewSubmitHandler(database *sqlx.DB, sugar *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		wordID, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_id"})
			return
		}

		var req reviewRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request"})
			return
		}

		userID, _ := c.Get("user_id")
		uid, _ := userID.(int)

		var word models.VocabWord
		if err := database.GetContext(c.Request.Context(), &word,
			"SELECT * FROM vocabulary WHERE id = ? AND user_id = ?", wordID, uid); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "word_not_found"})
			return
		}

		card := ai.SM2Card{
			Repetitions: word.ReviewCount,
			Interval:    word.Interval,
			EaseFactor:  word.EaseFactor,
		}

		result := ai.ProcessReview(card, req.Quality)

		nextReview := result.NextReview.Format(time.RFC3339)

		if err := db.UpdateReview(c.Request.Context(), database, wordID, uid,
			result.Repetitions, result.Interval, result.EaseFactor, nextReview); err != nil {
			sugar.Errorw("update review", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error"})
			return
		}

		sugar.Debugw("review submit",
			"word", word.Word,
			"quality", req.Quality,
			"interval", result.Interval,
			"ease_factor", result.EaseFactor,
			"next_review", nextReview,
		)

		c.JSON(http.StatusOK, reviewResponse{NextReview: nextReview})
	}
}

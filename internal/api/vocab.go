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

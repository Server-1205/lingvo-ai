package api

import (
	"encoding/csv"
	"encoding/json"
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
	WordID  int `json:"word_id" binding:"required"`
	Quality int `json:"quality" binding:"min=0,max=5"`
}

type reviewResponse struct {
	NextReview string `json:"next_review"`
}

func vocabListHandler(database *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")

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
		dueOnly := c.DefaultQuery("due_only", "false") == "true"

		where := "WHERE user_id = ?"
		args := []interface{}{uid}
		if dueOnly {
			where += " AND (next_review IS NULL OR next_review <= datetime('now'))"
		}

		var total int
		if err := database.GetContext(c.Request.Context(), &total,
			"SELECT COUNT(*) FROM vocabulary "+where, args...); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error"})
			return
		}

		query := "SELECT * FROM vocabulary " + where + " ORDER BY created_at DESC LIMIT ? OFFSET ?"
		queryArgs := append(args, perPage, (page-1)*perPage)
		var words []models.VocabWord
		if err := database.SelectContext(c.Request.Context(), &words, query, queryArgs...); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error"})
			return
		}
		if words == nil {
			words = []models.VocabWord{}
		}

		dueCount, _ := db.GetDueWordCount(c.Request.Context(), database, uid)

		c.JSON(http.StatusOK, models.VocabListResponse{
			Words:    words,
			Total:    total,
			DueCount: dueCount,
		})
	}
}

func vocabAddHandler(database *sqlx.DB, aiClient *ai.Client, sugar *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")

		if aiClient == nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "ai_service_unavailable"})
			return
		}

		var req struct {
			Word string `json:"word" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request"})
			return
		}

		userID, _ := c.Get("user_id")
		uid, _ := userID.(int)
		lang, _ := c.Get("lang")
		lng, _ := lang.(string)

		prompt := ai.BuildVocabPrompt(lng, req.Word)
		raw, err := aiClient.Generate(c.Request.Context(), prompt)
		if err != nil {
			sugar.Errorw("ai vocab lookup error", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ai_service_unavailable"})
			return
		}

		raw = cleanJSON(raw)
		var lookupResp models.VocabLookupResponse
		if err := json.Unmarshal([]byte(raw), &lookupResp); err != nil {
			sugar.Errorw("parse vocab response", "error", err, "raw", raw)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "parse_error"})
			return
		}

		example := ""
		if len(lookupResp.Examples) > 0 {
			example = lookupResp.Examples[0]
		}
		if err := db.AddVocabulary(c.Request.Context(), database, uid, req.Word, lookupResp.TranslationUz, example, lookupResp.Level); err != nil {
			sugar.Errorw("add vocab", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error"})
			return
		}

		sugar.Infow("vocab added via ai", "word", req.Word, "user_id", uid)

		c.JSON(http.StatusCreated, lookupResp)
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

func vocabExportHandler(database *sqlx.DB, sugar *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")
		c.Header("Content-Type", "text/csv")
		c.Header("Content-Disposition", `attachment; filename="vocabulary.csv"`)

		userID, _ := c.Get("user_id")
		uid, _ := userID.(int)

		var words []models.VocabWord
		if err := database.SelectContext(c.Request.Context(), &words,
			"SELECT * FROM vocabulary WHERE user_id = ? ORDER BY created_at DESC", uid); err != nil {
			sugar.Errorw("vocab export db error", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error"})
			return
		}

		writer := csv.NewWriter(c.Writer)
		writer.Write([]string{"word", "translation", "example", "level", "review_count", "next_review", "created_at"})

		for _, w := range words {
			nextReview := ""
			if w.NextReview != nil {
				nextReview = w.NextReview.Format("2006-01-02")
			}
			writer.Write([]string{
				w.Word,
				w.Translation,
				w.Example,
				w.Level,
				strconv.Itoa(w.ReviewCount),
				nextReview,
				w.CreatedAt.Format("2006-01-02"),
			})
		}

		writer.Flush()
		if err := writer.Error(); err != nil {
			sugar.Errorw("vocab export csv write error", "error", err)
		}

		sugar.Infow("vocab export", "user_id", uid, "count", len(words))
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
		c.Header("Cache-Control", "no-store")

		var req reviewRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request"})
			return
		}

		userID, _ := c.Get("user_id")
		uid, _ := userID.(int)

		var word models.VocabWord
		if err := database.GetContext(c.Request.Context(), &word,
			"SELECT * FROM vocabulary WHERE id = ? AND user_id = ?", req.WordID, uid); err != nil {
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

		if err := db.UpdateReview(c.Request.Context(), database, req.WordID, uid,
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

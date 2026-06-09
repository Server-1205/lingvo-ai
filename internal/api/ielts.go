package api

import (
	"encoding/json"
	"fmt"
	"math"
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

func ieltsWritingHandler(database *sqlx.DB, aiClient *ai.Client, sugar *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")

		isPremium, _ := c.Get("is_premium")
		if isPremium != true {
			c.JSON(http.StatusForbidden, gin.H{"error": "premium_required"})
			return
		}

		if aiClient == nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "ai_service_unavailable"})
			return
		}

		var req models.IeltsWritingRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request"})
			return
		}

		if req.Type != "task1" && req.Type != "task2" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_type"})
			return
		}

		lang, _ := c.Get("lang")
		lng, _ := lang.(string)

		prompt := ai.BuildIeltsWritingPrompt(req.Type, lng, req.UserText, req.TaskDescription)
		sugar.Debugw("ielts writing prompt built", "type", req.Type, "lang", lng)

		raw, err := aiClient.Generate(c.Request.Context(), prompt)
		if err != nil {
			sugar.Errorw("ielts writing ai error", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ai_service_unavailable"})
			return
		}

		raw = cleanJSON(raw)

		var resp models.IeltsWritingResponse
		if err := json.Unmarshal([]byte(raw), &resp); err != nil {
			sugar.Errorw("parse ielts writing response", "error", err, "raw", raw)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "parse_error"})
			return
		}

		details, _ := json.Marshal(resp.Criteria)
		userID, _ := c.Get("user_id")
		uid, _ := userID.(int)

		module := "writing_task1"
		if req.Type == "task2" {
			module = "writing_task2"
		}

		if saveErr := db.SaveIeltsScore(c.Request.Context(), database, uid, module, resp.BandScore, string(details), req.TaskDescription, req.UserText, resp.Feedback); saveErr != nil {
			sugar.Warnw("failed to save ielts score", "error", saveErr)
		}

		sugar.Debugw("ielts writing evaluated", "user_id", uid, "band_score", resp.BandScore)
		c.JSON(http.StatusOK, resp)
	}
}

func ieltsSpeakingGenerateHandler(database *sqlx.DB, aiClient *ai.Client, sugar *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")

		isPremium, _ := c.Get("is_premium")
		if isPremium != true {
			c.JSON(http.StatusForbidden, gin.H{"error": "premium_required"})
			return
		}

		if aiClient == nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "ai_service_unavailable"})
			return
		}

		partStr := c.DefaultQuery("part", "1")
		part, err := strconv.Atoi(partStr)
		if err != nil || part < 1 || part > 3 {
			part = 1
		}

		lang, _ := c.Get("lang")
		lng, _ := lang.(string)

		prompt := ai.BuildIeltsSpeakingPrompt(part, lng)
		sugar.Debugw("ielts speaking prompt built", "part", part, "lang", lng)

		raw, err := aiClient.Generate(c.Request.Context(), prompt)
		if err != nil {
			sugar.Errorw("ielts speaking ai error", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ai_service_unavailable"})
			return
		}

		raw = cleanJSON(raw)

		var resp models.IeltsSpeakingQuestionsResponse
		if err := json.Unmarshal([]byte(raw), &resp); err != nil {
			sugar.Errorw("parse ielts speaking response", "error", err, "raw", raw)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "parse_error"})
			return
		}

		resp.Part = part
		sugar.Debugw("ielts speaking questions generated", "part", part)
		c.JSON(http.StatusOK, resp)
	}
}

func ieltsSpeakingEvaluateHandler(database *sqlx.DB, aiClient *ai.Client, sugar *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")

		isPremium, _ := c.Get("is_premium")
		if isPremium != true {
			c.JSON(http.StatusForbidden, gin.H{"error": "premium_required"})
			return
		}

		if aiClient == nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "ai_service_unavailable"})
			return
		}

		var req models.IeltsSpeakingRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request"})
			return
		}

		if req.Part < 1 || req.Part > 3 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_part"})
			return
		}

		lang, _ := c.Get("lang")
		lng, _ := lang.(string)

		prompt := ai.BuildIeltsSpeakingEvaluatePrompt(req.Part, lng, req.Question, req.UserResponse)
		sugar.Debugw("ielts speaking evaluate prompt built", "part", req.Part)

		raw, err := aiClient.Generate(c.Request.Context(), prompt)
		if err != nil {
			sugar.Errorw("ielts speaking evaluate ai error", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ai_service_unavailable"})
			return
		}

		raw = cleanJSON(raw)

		var resp models.IeltsSpeakingResponse
		if err := json.Unmarshal([]byte(raw), &resp); err != nil {
			sugar.Errorw("parse ielts speaking evaluate response", "error", err, "raw", raw)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "parse_error"})
			return
		}

		details, _ := json.Marshal(resp.Criteria)
		userID, _ := c.Get("user_id")
		uid, _ := userID.(int)

		if saveErr := db.SaveIeltsScore(c.Request.Context(), database, uid, "speaking", resp.BandScore, string(details), req.Question, req.UserResponse, resp.Feedback); saveErr != nil {
			sugar.Warnw("failed to save ielts score", "error", saveErr)
		}

		sugar.Debugw("ielts speaking evaluated", "user_id", uid, "band_score", resp.BandScore)
		c.JSON(http.StatusOK, resp)
	}
}

func ieltsReadingGenerateHandler(database *sqlx.DB, aiClient *ai.Client, sugar *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")

		isPremium, _ := c.Get("is_premium")
		if isPremium != true {
			c.JSON(http.StatusForbidden, gin.H{"error": "premium_required"})
			return
		}

		if aiClient == nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "ai_service_unavailable"})
			return
		}

		lang, _ := c.Get("lang")
		lng, _ := lang.(string)

		prompt := ai.BuildIeltsReadingPrompt(lng)
		sugar.Debugw("ielts reading prompt built", "lang", lng)

		raw, err := aiClient.Generate(c.Request.Context(), prompt)
		if err != nil {
			sugar.Errorw("ielts reading ai error", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ai_service_unavailable"})
			return
		}

		raw = cleanJSON(raw)

		var resp models.IeltsReadingPassage
		if err := json.Unmarshal([]byte(raw), &resp); err != nil {
			sugar.Errorw("parse ielts reading response", "error", err, "raw", raw)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "parse_error"})
			return
		}

		sugar.Debugw("ielts reading passage generated", "word_count", resp.WordCount)
		c.JSON(http.StatusOK, resp)
	}
}

func ieltsReadingEvaluateHandler(database *sqlx.DB, aiClient *ai.Client, sugar *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")

		isPremium, _ := c.Get("is_premium")
		if isPremium != true {
			c.JSON(http.StatusForbidden, gin.H{"error": "premium_required"})
			return
		}

		if aiClient == nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "ai_service_unavailable"})
			return
		}

		var req models.IeltsReadingSubmitRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request"})
			return
		}

		if len(req.UserAnswers) != len(req.Questions) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "answers_mismatch"})
			return
		}

		lang, _ := c.Get("lang")
		lng, _ := lang.(string)

		questionsJSON, _ := json.Marshal(req.Questions)
		answersJSON, _ := json.Marshal(req.UserAnswers)

		prompt := ai.BuildIeltsReadingEvaluatePrompt(lng, req.Passage, string(questionsJSON), string(answersJSON))
		sugar.Debugw("ielts reading evaluate prompt built")

		raw, err := aiClient.Generate(c.Request.Context(), prompt)
		if err != nil {
			sugar.Errorw("ielts reading evaluate ai error", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ai_service_unavailable"})
			return
		}

		raw = cleanJSON(raw)

		var resp models.IeltsReadingResult
		if err := json.Unmarshal([]byte(raw), &resp); err != nil {
			sugar.Errorw("parse ielts reading evaluate response", "error", err, "raw", raw)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "parse_error"})
			return
		}

		details := fmt.Sprintf(`{"correct":%d,"total":%d}`, resp.CorrectAnswers, resp.TotalQuestions)
		userID, _ := c.Get("user_id")
		uid, _ := userID.(int)

		feedbackJSON, _ := json.Marshal(resp)

		if saveErr := db.SaveIeltsScore(c.Request.Context(), database, uid, "reading", resp.BandScore, details, req.Passage, string(answersJSON), string(feedbackJSON)); saveErr != nil {
			sugar.Warnw("failed to save ielts score", "error", saveErr)
		}

		sugar.Debugw("ielts reading evaluated", "user_id", uid, "band_score", resp.BandScore)
		c.JSON(http.StatusOK, resp)
	}
}

func ieltsScoresHandler(database *sqlx.DB, sugar *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")

		isPremium, _ := c.Get("is_premium")
		if isPremium != true {
			c.JSON(http.StatusForbidden, gin.H{"error": "premium_required"})
			return
		}

		userID, _ := c.Get("user_id")
		uid, _ := userID.(int)

		module := c.Query("module")
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))
		if page < 1 {
			page = 1
		}
		if perPage < 1 || perPage > 50 {
			perPage = 10
		}
		offset := (page - 1) * perPage

		entries, total, err := db.GetIeltsScores(c.Request.Context(), database, uid, module, perPage, offset)
		if err != nil {
			sugar.Errorw("ielts scores query failed", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error"})
			return
		}

		stats, _ := db.GetIeltsScoreStats(c.Request.Context(), database, uid)

		sugar.Debugw("ielts scores requested", "user_id", uid, "total", total)
		c.JSON(http.StatusOK, models.IeltsScoresResponse{
			Entries: entries,
			Total:   total,
			Stats:   stats,
		})
	}
}

func readingBandScore(correct, total int) float64 {
	if total == 0 {
		return 0
	}
	ratio := float64(correct) / float64(total)
	switch {
	case ratio >= 0.95:
		return 9.0
	case ratio >= 0.85:
		return 8.0
	case ratio >= 0.75:
		return 7.0
	case ratio >= 0.65:
		return 6.5
	case ratio >= 0.55:
		return 6.0
	case ratio >= 0.45:
		return 5.5
	case ratio >= 0.35:
		return 5.0
	default:
		return 4.0
	}
}

func roundBandScore(raw float64) float64 {
	return math.Round(raw*2) / 2
}

func init() {
	_ = time.Now
	_ = roundBandScore
	_ = readingBandScore
}

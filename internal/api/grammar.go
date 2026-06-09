package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/lingvo-ai/lingvo/internal/ai"
	"github.com/lingvo-ai/lingvo/internal/models"
)

func grammarHandler(database *sqlx.DB, aiClient *ai.Client, sugar *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")

		var req models.GrammarRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request"})
			return
		}

		lang, _ := c.Get("lang")
		level, _ := c.Get("level")
		lvl, _ := level.(string)
		lng, _ := lang.(string)

		prompt := ai.BuildGrammarCheckPrompt(lvl, lng, req.Text)

		raw, err := aiClient.GenerateLite(c.Request.Context(), prompt)
		if err != nil {
			sugar.Errorw("ai grammar error", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ai_service_unavailable"})
			return
		}

		raw = cleanJSON(raw)

		var grammarResp models.AIResponse
		if err := json.Unmarshal([]byte(raw), &grammarResp); err != nil {
			sugar.Errorw("parse grammar response", "error", err, "raw", raw)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "parse_error"})
			return
		}

		var corrections []models.Correction
		if grammarResp.Corrections != nil {
			for _, cr := range grammarResp.Corrections {
				corrections = append(corrections, models.Correction{
					Original:      cr.Original,
					Corrected:     cr.Corrected,
					ExplanationUz: cr.ExplanationUz,
					ExplanationRu: cr.ExplanationRu,
					Type:          cr.Type,
				})
			}
		}

		c.JSON(http.StatusOK, models.GrammarResponse{Corrections: corrections})
	}
}

func vocabLookupHandler(aiClient *ai.Client, sugar *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.VocabLookupRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request"})
			return
		}

		lng := req.Lang
		if lng == "" {
			if ctxLang, ok := c.Get("lang"); ok {
				lng, _ = ctxLang.(string)
			}
		}

		if lng == "uz" && hasCyrillic(req.Word) {
			sugar.Warnw("[vocab] cyrillic word rejected in uz mode", "word", req.Word)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_language"})
			return
		}

		prompt := ai.BuildVocabPrompt(lng, req.Word)

		raw, err := aiClient.GenerateLite(c.Request.Context(), prompt)
		if err != nil {
			sugar.Errorw("ai vocab lookup error", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ai_service_unavailable"})
			return
		}

		raw = cleanJSON(raw)

		var resp models.VocabLookupResponse
		if err := json.Unmarshal([]byte(raw), &resp); err != nil {
			sugar.Errorw("parse vocab response", "error", err, "raw", raw)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "parse_error"})
			return
		}

		if resp.Error != "" {
			sugar.Warnw("[vocab] AI rejected word", "word", req.Word, "reason", resp.Error)
			c.JSON(http.StatusBadRequest, gin.H{"error": resp.Error})
			return
		}

		c.JSON(http.StatusOK, resp)
	}
}

func cleanJSON(raw string) string {
	raw = strings.TrimSpace(raw)
	raw = strings.TrimPrefix(raw, "```json")
	raw = strings.TrimPrefix(raw, "```")
	raw = strings.TrimSuffix(raw, "```")
	return strings.TrimSpace(raw)
}

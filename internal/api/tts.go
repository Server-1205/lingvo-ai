package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/lingvo-ai/lingvo/internal/tts"
)

const maxTTSLength = 500

func ttsHandler(ttsClient tts.Synthesizer, sugar *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")

		text := c.Query("text")
		lang := c.DefaultQuery("lang", "uz")

		if text == "" {
			sugar.Warnw("[tts] missing text parameter")
			c.JSON(http.StatusBadRequest, gin.H{"error": "text_required"})
			return
		}

		if len(text) > maxTTSLength {
			sugar.Warnw("[tts] text too long", "len", len(text))
			c.JSON(http.StatusBadRequest, gin.H{"error": "text_too_long"})
			return
		}

		if lang != "uz" && lang != "ru" {
			lang = "uz"
		}

		sugar.Debugw("[tts] request", "text", text, "lang", lang)

		audio, err := ttsClient.Synthesize(c.Request.Context(), text, lang)
		if err != nil {
			sugar.Errorw("[tts] synthesis error", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "tts_failed"})
			return
		}

		c.Data(http.StatusOK, "audio/mpeg", audio)
	}
}

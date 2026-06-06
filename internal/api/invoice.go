package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/lingvo-ai/lingvo/internal/models"
)

var planPrices = map[string]int{
	"weekly":  300,
	"monthly": 800,
}

var planTitles = map[string]map[string]string{
	"weekly": {
		"uz": "Lingvo AI — Ҳафталик обуна",
		"ru": "Lingvo AI — Недельная подписка",
	},
	"monthly": {
		"uz": "Lingvo AI — Ойлик обуна",
		"ru": "Lingvo AI — Месячная подписка",
	},
}

var planDescs = map[string]map[string]string{
	"weekly": {
		"uz": "Чексиз AI хабарлар. 7 кун амал қилади.",
		"ru": "Неограниченные AI-сообщения. Действует 7 дней.",
	},
	"monthly": {
		"uz": "Чексиз AI хабарлар. 30 кун амал қилади.",
		"ru": "Неограниченные AI-сообщения. Действует 30 дней.",
	},
}

func invoiceHandler(botToken string, sugar *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.InvoiceRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "invalid_request"})
			return
		}

		price, ok := planPrices[req.Plan]
		if !ok {
			c.JSON(400, gin.H{"error": "invalid_plan"})
			return
		}

		lang, _ := c.Get("lang")
		lng, _ := lang.(string)

		title := planTitles[req.Plan][lng]
		desc := planDescs[req.Plan][lng]
		if title == "" {
			title = planTitles[req.Plan]["uz"]
			desc = planDescs[req.Plan]["uz"]
		}

		payload := fmt.Sprintf("%s_%d", req.Plan, c.GetInt("user_id"))

		link, err := createTelegramInvoiceLink(botToken, title, desc, payload, price)
		if err != nil {
			sugar.Errorw("create invoice link", "error", err)
			c.JSON(500, gin.H{"error": "internal_error"})
			return
		}

		c.JSON(200, models.InvoiceResponse{
			InvoiceLink: link,
			Stars:       price,
		})
	}
}

func createTelegramInvoiceLink(botToken, title, description, payload string, amount int) (string, error) {
	body := map[string]interface{}{
		"title":              title,
		"description":        description,
		"payload":            payload,
		"currency":           "XTR",
		"prices": []map[string]interface{}{
			{"label": title, "amount": amount},
		},
	}

	b, _ := json.Marshal(body)
	url := fmt.Sprintf("https://api.telegram.org/bot%s/createInvoiceLink", botToken)

	resp, err := http.Post(url, "application/json", bytes.NewReader(b))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Ok          bool   `json:"ok"`
		Result      string `json:"result"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if !result.Ok {
		return "", fmt.Errorf("telegram api error: %s", result.Description)
	}

	return result.Result, nil
}

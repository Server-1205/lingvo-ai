package ai

import (
	"encoding/json"
	"strings"

	"github.com/lingvo-ai/lingvo/internal/models"
)

type StreamEventType string

const (
	EventToken       StreamEventType = "token"
	EventResult      StreamEventType = "result"
)

type StreamEvent struct {
	Type StreamEventType `json:"type"`
	Data json.RawMessage `json:"data"`
}

func ParseAIResponse(raw string) (*models.AIResponse, error) {
	cleaned := strings.TrimSpace(raw)
	cleaned = strings.TrimPrefix(cleaned, "```json")
	cleaned = strings.TrimPrefix(cleaned, "```")
	cleaned = strings.TrimSuffix(cleaned, "```")
	cleaned = strings.TrimSpace(cleaned)

	var resp models.AIResponse
	if err := json.Unmarshal([]byte(cleaned), &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

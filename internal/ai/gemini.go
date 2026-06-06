package ai

import (
	"context"
	"encoding/json"

	"github.com/google/generative-ai-go/genai"
	"github.com/lingvo-ai/lingvo/internal/models"
	"google.golang.org/api/option"
)

type Client struct {
	genClient *genai.Client
	model     *genai.GenerativeModel
}

func NewClient(ctx context.Context, apiKey string) (*Client, error) {
	genClient, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}

	model := genClient.GenerativeModel("gemini-2.0-flash")
	model.SetTemperature(0.4)

	return &Client{
		genClient: genClient,
		model:     model,
	}, nil
}

func (c *Client) Close() error {
	return c.genClient.Close()
}

func (c *Client) Chat(ctx context.Context, prompt string) (*models.AIResponse, error) {
	resp, err := c.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, err
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return nil, nil
	}

	text := ""
	for _, part := range resp.Candidates[0].Content.Parts {
		if t, ok := part.(genai.Text); ok {
			text += string(t)
		}
	}

	if text == "" {
		return nil, nil
	}

	var aiResp models.AIResponse
	if err := json.Unmarshal([]byte(text), &aiResp); err != nil {
		return nil, err
	}

	if aiResp.Reply == "" {
		aiResp.Reply = text
	}

	return &aiResp, nil
}

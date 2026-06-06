package ai

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"github.com/lingvo-ai/lingvo/internal/models"
	"google.golang.org/api/iterator"
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

func (c *Client) Generate(ctx context.Context, prompt string) (string, error) {
	resp, err := c.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", err
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", nil
	}

	text := ""
	for _, part := range resp.Candidates[0].Content.Parts {
		if t, ok := part.(genai.Text); ok {
			text += string(t)
		}
	}

	return text, nil
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

func (c *Client) ChatStream(ctx context.Context, prompt string) <-chan StreamEvent {
	ch := make(chan StreamEvent)

	go func() {
		defer close(ch)

		iter := c.model.GenerateContentStream(ctx, genai.Text(prompt))

		var fullText strings.Builder

		for {
			chunk, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				errData, _ := json.Marshal(err.Error())
				ch <- StreamEvent{Type: EventToken, Data: errData}
				return
			}

			if len(chunk.Candidates) == 0 || len(chunk.Candidates[0].Content.Parts) == 0 {
				continue
			}

			for _, part := range chunk.Candidates[0].Content.Parts {
				if t, ok := part.(genai.Text); ok {
					text := string(t)
					fullText.WriteString(text)
					data, _ := json.Marshal(text)
					ch <- StreamEvent{Type: EventToken, Data: data}
				}
			}
		}

		if fullText.Len() == 0 {
			return
		}

		aiResp, err := ParseAIResponse(fullText.String())
		if err != nil {
			return
		}

		if aiResp.Reply == "" {
			aiResp.Reply = fullText.String()
		}

		resultData, err := json.Marshal(aiResp)
		if err != nil {
			return
		}
		ch <- StreamEvent{Type: EventResult, Data: resultData}
	}()

	return ch
}

package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"github.com/lingvo-ai/lingvo/internal/models"
	"go.uber.org/zap"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type Client struct {
	genClient     *genai.Client
	model         *genai.GenerativeModel
	liteModel     *genai.GenerativeModel
	fallback      *openAIClient
	sugar         *zap.SugaredLogger
	queue         *Queue
	queueEnabled  bool
}

func NewClient(ctx context.Context, apiKey string, openAIKey, openAIBaseURL, openAIModel string, sugar *zap.SugaredLogger) (*Client, error) {
	var fb *openAIClient
	if openAIKey != "" {
		fb = newOpenAIClient(openAIKey, openAIBaseURL, openAIModel)
		sugar.Infow("primary AI configured", "model", openAIModel, "base_url", openAIBaseURL)
	} else {
		sugar.Warn("OPENAI_API_KEY not set, OpenAI-compatible AI disabled")
	}

	c := &Client{
		fallback: fb,
		sugar:    sugar,
	}

	if apiKey != "" {
		genClient, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
		if err != nil {
			sugar.Warnw("failed to init Gemini client, fallback AI only", "error", err)
		} else {
			model := genClient.GenerativeModel("gemini-2.0-flash")
			model.SetTemperature(0.3)

			liteModel := genClient.GenerativeModel("gemini-2.0-flash-lite")
			liteModel.SetTemperature(0.2)

			c.genClient = genClient
			c.model = model
			c.liteModel = liteModel
			sugar.Info("Gemini fallback configured")
		}
	} else {
		sugar.Info("GEMINI_API_KEY not set, running with primary AI only")
	}

	if c.fallback == nil && c.model == nil {
		return nil, fmt.Errorf("no AI provider configured")
	}

	return c, nil
}

func (c *Client) Close() error {
	if c.queue != nil {
		c.queue.StopWorker()
	}
	if c.genClient != nil {
		return c.genClient.Close()
	}
	return nil
}

func (c *Client) EnableQueue(ctx context.Context, sugar *zap.SugaredLogger) {
	c.queue = NewQueue()
	c.queue.StartWorker(ctx, c, sugar)
	c.queueEnabled = true
}

func (c *Client) EnqueueAI(req *AIRequest) {
	if c.queue != nil {
		c.queue.Enqueue(req)
	}
}

func (c *Client) IsQueueEnabled() bool {
	return c.queueEnabled
}

func (c *Client) GenerateLite(ctx context.Context, prompt string) (string, error) {
	if c.fallback != nil && c.fallback.IsConfigured() {
		text, err := c.fallback.Generate(ctx, prompt)
		if err == nil {
			return text, nil
		}
		c.sugar.Warnw("primary AI (openai-compatible) lite failed, trying Gemini", "error", err)
	}

	if c.liteModel != nil {
		text, err := c.generateModel(ctx, c.liteModel, prompt)
		if err == nil {
			return text, nil
		}
		return "", err
	}

	return "", fmt.Errorf("no AI provider available")
}

func (c *Client) Generate(ctx context.Context, prompt string) (string, error) {
	if c.fallback != nil && c.fallback.IsConfigured() {
		text, err := c.fallback.Generate(ctx, prompt)
		if err == nil {
			return text, nil
		}
		c.sugar.Warnw("primary AI (openai-compatible) failed, trying Gemini", "error", err)
	}

	if c.model != nil {
		text, err := c.generateModel(ctx, c.model, prompt)
		if err == nil {
			return text, nil
		}
		return "", err
	}

	return "", fmt.Errorf("no AI provider available")
}

func (c *Client) generateModel(ctx context.Context, model *genai.GenerativeModel, prompt string) (string, error) {
	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
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
	if c.fallback != nil && c.fallback.IsConfigured() {
		fallbackText, fbErr := c.fallback.Generate(ctx, prompt)
		if fbErr == nil && fallbackText != "" {
			return c.parseAIResponse(fallbackText)
		}
		c.sugar.Warnw("primary AI (openai-compatible) chat failed, trying Gemini", "error", fbErr)
	}

	if c.model != nil {
		text, err := c.chatWithModel(ctx, c.model, prompt)
		if err == nil && text != "" {
			return c.parseAIResponse(text)
		}
		if text != "" {
			return c.parseAIResponse(text)
		}
		return nil, err
	}

	return nil, fmt.Errorf("no AI provider available")
}

func (c *Client) chatWithModel(ctx context.Context, model *genai.GenerativeModel, prompt string) (string, error) {
	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
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

func (c *Client) parseAIResponse(text string) (*models.AIResponse, error) {
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

		// Try OpenAI-compatible streaming first
		if c.fallback != nil && c.fallback.IsConfigured() {
			streamCh, err := c.fallback.GenerateStream(ctx, prompt)
			if err == nil {
				for evt := range streamCh {
					ch <- evt
				}
				return
			}
			c.sugar.Warnw("primary AI (openai-compatible) stream failed, trying Gemini", "error", err)
		}

		// Fallback to Gemini streaming
		if c.model == nil {
			errData, _ := json.Marshal("ai_service_unavailable")
			ch <- StreamEvent{Type: EventToken, Data: errData}
			return
		}

		iter := c.model.GenerateContentStream(ctx, genai.Text(prompt))
		var fullText strings.Builder

		for {
			chunk, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				c.sugar.Errorw("gemini stream also failed", "error", err)
				errData, _ := json.Marshal("ai_service_unavailable")
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

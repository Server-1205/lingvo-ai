package ai

import (
	"context"
	"encoding/json"
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
	genClient, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}

	model := genClient.GenerativeModel("gemini-2.0-flash")
	model.SetTemperature(0.4)

	liteModel := genClient.GenerativeModel("gemini-2.0-flash-lite")
	liteModel.SetTemperature(0.4)

	var fb *openAIClient
	if openAIKey != "" {
		fb = newOpenAIClient(openAIKey, openAIBaseURL, openAIModel)
		sugar.Infow("openai-compatible fallback configured", "model", openAIModel, "base_url", openAIBaseURL)
	} else {
		sugar.Warn("OPENAI_API_KEY not set, AI fallback disabled")
	}

	return &Client{
		genClient: genClient,
		model:     model,
		liteModel: liteModel,
		fallback:  fb,
		sugar:     sugar,
	}, nil
}

func (c *Client) Close() error {
	if c.queue != nil {
		c.queue.StopWorker()
	}
	return c.genClient.Close()
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
	text, err := c.generateModel(ctx, c.liteModel, prompt)
	if err == nil {
		return text, nil
	}

	c.sugar.Warnw("gemini lite generate failed, trying fallback", "error", err)
	if c.fallback != nil && c.fallback.IsConfigured() {
		fallbackText, fbErr := c.fallback.Generate(ctx, prompt)
		if fbErr == nil {
			c.sugar.Warn("fallback AI used for lite request")
			return fallbackText, nil
		}
		c.sugar.Errorw("fallback also failed", "error", fbErr)
	}

	return "", err
}

func (c *Client) Generate(ctx context.Context, prompt string) (string, error) {
	text, err := c.generateModel(ctx, c.model, prompt)
	if err == nil {
		return text, nil
	}

	c.sugar.Warnw("gemini generate failed, trying fallback", "error", err)
	if c.fallback != nil && c.fallback.IsConfigured() {
		fallbackText, fbErr := c.fallback.Generate(ctx, prompt)
		if fbErr == nil {
			c.sugar.Warn("fallback AI used")
			return fallbackText, nil
		}
		c.sugar.Errorw("fallback also failed", "error", fbErr)
	}

	return "", err
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
	text, err := c.chatWithModel(ctx, c.model, prompt)
	if err == nil && text != "" {
		return c.parseAIResponse(text)
	}

	if err != nil {
		c.sugar.Warnw("gemini chat failed, trying fallback", "error", err)
	} else {
		c.sugar.Warn("gemini chat returned empty, trying fallback")
	}

	if c.fallback != nil && c.fallback.IsConfigured() {
		fallbackText, fbErr := c.fallback.Generate(ctx, prompt)
		if fbErr == nil && fallbackText != "" {
			c.sugar.Warn("fallback AI used for chat")
			return c.parseAIResponse(fallbackText)
		}
		c.sugar.Errorw("fallback also failed for chat", "error", fbErr)
	}

	if text != "" {
		return c.parseAIResponse(text)
	}
	return nil, err
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

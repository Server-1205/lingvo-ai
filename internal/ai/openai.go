package ai

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type openAIClient struct {
	baseURL string
	apiKey  string
	model   string
	client  *http.Client
}

type openAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAIRequest struct {
	Model       string          `json:"model"`
	Messages    []openAIMessage `json:"messages"`
	Temperature float64         `json:"temperature"`
	Stream      bool            `json:"stream,omitempty"`
}

type openAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

type openAIStreamChunk struct {
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
		FinishReason *string `json:"finish_reason"`
	} `json:"choices"`
}

func newOpenAIClient(apiKey, baseURL, model string) *openAIClient {
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}
	if model == "" {
		model = "gpt-4o-mini"
	}
	return &openAIClient{
		baseURL: baseURL,
		apiKey:  apiKey,
		model:   model,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (o *openAIClient) IsConfigured() bool {
	return o.apiKey != ""
}

func (o *openAIClient) doRequest(ctx context.Context, body openAIRequest) (string, error) {
	payload, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("openai marshal: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, o.baseURL+"/chat/completions", bytes.NewReader(payload))
	if err != nil {
		return "", fmt.Errorf("openai request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+o.apiKey)

	resp, err := o.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("openai do: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("openai read: %w", err)
	}

	var result openAIResponse
	if err := json.Unmarshal(raw, &result); err != nil {
		return "", fmt.Errorf("openai unmarshal: %w", err)
	}

	if result.Error != nil {
		return "", fmt.Errorf("openai api error: %s", result.Error.Message)
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("openai: no choices")
	}

	return result.Choices[0].Message.Content, nil
}

func (o *openAIClient) Generate(ctx context.Context, prompt string) (string, error) {
	body := openAIRequest{
		Model: o.model,
		Messages: []openAIMessage{
			{Role: "user", Content: prompt},
		},
		Temperature: 0.3,
	}
	return o.doRequest(ctx, body)
}

func (o *openAIClient) GenerateStream(ctx context.Context, prompt string) (<-chan StreamEvent, error) {
	body := openAIRequest{
		Model: o.model,
		Messages: []openAIMessage{
			{Role: "user", Content: prompt},
		},
		Temperature: 0.3,
		Stream:      true,
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("openai stream marshal: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, o.baseURL+"/chat/completions", bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("openai stream request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+o.apiKey)
	req.Header.Set("Accept", "text/event-stream")

	resp, err := o.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("openai stream do: %w", err)
	}

	ch := make(chan StreamEvent)

	go func() {
		defer close(ch)
		defer resp.Body.Close()

		scanner := bufio.NewScanner(resp.Body)
		var fullText strings.Builder

		for scanner.Scan() {
			line := scanner.Text()

			if !strings.HasPrefix(line, "data: ") {
				continue
			}

			data := strings.TrimPrefix(line, "data: ")
			if data == "[DONE]" {
				break
			}

			var chunk openAIStreamChunk
			if err := json.Unmarshal([]byte(data), &chunk); err != nil {
				continue
			}

			if len(chunk.Choices) == 0 {
				continue
			}

			content := chunk.Choices[0].Delta.Content
			if content == "" {
				continue
			}

			fullText.WriteString(content)
			tokenData, _ := json.Marshal(content)
			ch <- StreamEvent{Type: EventToken, Data: tokenData}
		}

		if err := scanner.Err(); err != nil {
			return
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

	return ch, nil
}

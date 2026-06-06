---
name: gemini-prompt
description: Use when designing or modifying AI prompts for Gemini integration, parsing AI responses (JSON format), or implementing the chat/correction pipeline for Lingvo AI. NOT for general API routing or database queries.
---

# Gemini AI Integration

## Client Setup (`internal/ai/client.go`)

```go
import "github.com/google/generative-ai-go/genai"

func NewClient(ctx context.Context, apiKey string) (*genai.Client, error) {
    client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
    if err != nil {
        return nil, fmt.Errorf("genai client: %w", err)
    }
    return client, nil
}
```

## Prompt Format (`internal/ai/prompts.go`)

```go
const systemPrompt = `You are an English tutor for a Uzbek-speaking student (level: %s).
Respond naturally to the student's message, then provide corrections.

You MUST respond with valid JSON only (no markdown, no code fences):
{
  "reply": "Your natural response in English...",
  "corrections": [
    {
      "original": "student's incorrect phrase",
      "corrected": "corrected version",
      "explanation": "explanation in Russian or Uzbek why this was wrong"
    }
  ]
}

If no corrections needed, corrections array must be empty.
Keep reply natural and encouraging. Level-appropriate vocabulary.`
```

## Parsing Response

```go
type ChatResponse struct {
    Reply       string       `json:"reply"`
    Corrections []Correction `json:"corrections"`
}

type Correction struct {
    Original    string `json:"original"`
    Corrected   string `json:"corrected"`
    Explanation string `json:"explanation"`
}

func parseResponse(text string) (*ChatResponse, error) {
    text = strings.TrimSpace(text)
    text = strings.TrimPrefix(text, "```json")
    text = strings.TrimSuffix(text, "```")
    text = strings.TrimSpace(text)

    var resp ChatResponse
    if err := json.Unmarshal([]byte(text), &resp); err != nil {
        return nil, fmt.Errorf("parse AI response: %w", err)
    }
    return &resp, nil
}
```

## Modes
- `chat` — conversational with corrections (above prompt)
- `grammar` — grammar-only check (returns corrections only, no reply)
- `quiz` — generate quiz questions based on vocabulary
- `level` — assess CEFR level from conversation sample

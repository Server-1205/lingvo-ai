package ai

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/lingvo-ai/lingvo/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func newTestSugar() *zap.SugaredLogger {
	return zap.NewNop().Sugar()
}

func TestChatStream_ReturnsChannel(t *testing.T) {
	key := os.Getenv("GEMINI_API_KEY")
	if key == "" {
		t.Skip("GEMINI_API_KEY not set")
	}
	sugar := newTestSugar()
	client, err := NewClient(context.Background(), key, "", "", "", sugar)
	require.NoError(t, err)
	defer client.Close()

	ch := client.ChatStream(context.Background(), "say hello in one word")
	require.NotNil(t, ch)

	var tokens int
	for evt := range ch {
		switch evt.Type {
		case EventToken:
			tokens++
		case EventResult:
			var result models.AIResponse
			err := json.Unmarshal(evt.Data, &result)
			require.NoError(t, err)
			assert.NotEmpty(t, result.Reply)
		}
	}
	assert.Greater(t, tokens, 0, "expected at least one token event")
}

func TestChatStream_EmitsTokenAndResult(t *testing.T) {
	key := os.Getenv("GEMINI_API_KEY")
	if key == "" {
		t.Skip("GEMINI_API_KEY not set")
	}
	sugar := newTestSugar()
	client, err := NewClient(context.Background(), key, "", "", "", sugar)
	require.NoError(t, err)
	defer client.Close()

	ch := client.ChatStream(context.Background(), "respond only with: {\"reply\":\"test\",\"corrections\":[]}")

	var gotToken, gotResult bool
	for evt := range ch {
		switch evt.Type {
		case EventToken:
			gotToken = true
		case EventResult:
			gotResult = true
			var result models.AIResponse
			err := json.Unmarshal(evt.Data, &result)
			require.NoError(t, err)
		}
	}
	assert.True(t, gotToken, "expected token events")
	assert.True(t, gotResult, "expected result event")
}

func TestChatStream_ContextCancellation(t *testing.T) {
	key := os.Getenv("GEMINI_API_KEY")
	if key == "" {
		t.Skip("GEMINI_API_KEY not set")
	}
	sugar := newTestSugar()
	client, err := NewClient(context.Background(), key, "", "", "", sugar)
	require.NoError(t, err)
	defer client.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	ch := client.ChatStream(ctx, "say hello in one word")
	require.NotNil(t, ch)

	for range ch {
	}
}

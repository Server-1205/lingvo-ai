package api

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type mockTTSClient struct {
	audio []byte
	err   error
}

func (m *mockTTSClient) Synthesize(_ context.Context, text, lang string) ([]byte, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.audio, nil
}

func (m *mockTTSClient) IsAvailable() bool {
	return true
}

func TestTTSHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger, _ := zap.NewDevelopment()
	sugar := logger.Sugar()

	mock := &mockTTSClient{audio: []byte{0xFF, 0xF3, 0xE4}}
	handler := ttsHandler(mock, sugar)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/tts?text=hello", nil)

	handler(c)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "audio/mpeg", w.Header().Get("Content-Type"))
	require.Equal(t, "no-store", w.Header().Get("Cache-Control"))
	require.Equal(t, []byte{0xFF, 0xF3, 0xE4}, w.Body.Bytes())
}

func TestTTSHandler_MissingText(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger, _ := zap.NewDevelopment()
	sugar := logger.Sugar()

	handler := ttsHandler(&mockTTSClient{}, sugar)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/tts", nil)

	handler(c)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestTTSHandler_TextTooLong(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger, _ := zap.NewDevelopment()
	sugar := logger.Sugar()

	handler := ttsHandler(&mockTTSClient{}, sugar)

	longText := make([]byte, 501)
	for i := range longText {
		longText[i] = 'a'
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/tts?text="+string(longText), nil)

	handler(c)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestTTSHandler_Failure(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger, _ := zap.NewDevelopment()
	sugar := logger.Sugar()

	mock := &mockTTSClient{err: errors.New("synthesis failed")}
	handler := ttsHandler(mock, sugar)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/tts?text=hello", nil)

	handler(c)

	require.Equal(t, http.StatusInternalServerError, w.Code)
}

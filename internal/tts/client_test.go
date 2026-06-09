package tts

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestSynthesize_EmptyText(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	client := NewClient("uz-UZ-MadinaNeural", "ru-RU-SvetlanaNeural", logger.Sugar())
	_, err := client.Synthesize(context.Background(), "", "uz")
	require.ErrorContains(t, err, "empty text")
}

func TestVoiceSelection_Uz(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	client := NewClient("uz-voice", "ru-voice", logger.Sugar())
	require.Equal(t, "uz-voice", client.voiceUz)
	require.Equal(t, "ru-voice", client.voiceRu)
}

func TestDefaultVoices(t *testing.T) {
	require.Equal(t, "en-US-JennyNeural", DefaultVoices.Uz)
	require.Equal(t, "en-US-JennyNeural", DefaultVoices.Ru)
}

package tts

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"go.uber.org/zap"
)

type Synthesizer interface {
	Synthesize(ctx context.Context, text, lang string) ([]byte, error)
	IsAvailable() bool
}

type Client struct {
	voiceUz string
	voiceRu string
	sugar   *zap.SugaredLogger
}

func NewClient(voiceUz, voiceRu string, sugar *zap.SugaredLogger) *Client {
	return &Client{
		voiceUz: voiceUz,
		voiceRu: voiceRu,
		sugar:   sugar,
	}
}

func (c *Client) Synthesize(ctx context.Context, text, lang string) ([]byte, error) {
	voice := c.voiceUz
	if lang == "ru" {
		voice = c.voiceRu
	}

	c.sugar.Debugw("[tts] starting synthesis", "text", text, "lang", lang, "voice", voice, "len", len(text))

	if text == "" {
		return nil, fmt.Errorf("tts: empty text")
	}

	tmpFile, err := os.CreateTemp("", "tts-*.mp3")
	if err != nil {
		return nil, fmt.Errorf("tts: create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()
	tmpFile.Close()
	defer os.Remove(tmpPath)

	start := time.Now()
	cmd := exec.CommandContext(ctx, "edge-tts",
		"--voice", voice,
		"--text", text,
		"--write-media", tmpPath,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		c.sugar.Errorw("[tts] synthesis failed", "error", err, "output", string(output), "text", text)
		return nil, fmt.Errorf("tts: edge-tts failed: %w", err)
	}

	audio, err := os.ReadFile(tmpPath)
	if err != nil {
		return nil, fmt.Errorf("tts: read audio file: %w", err)
	}

	elapsed := time.Since(start)
	c.sugar.Debugw("[tts] synthesized", "bytes", len(audio), "duration", elapsed, "lang", lang)

	return audio, nil
}

func (c *Client) IsAvailable() bool {
	_, err := exec.LookPath("edge-tts")
	if err != nil {
		c.sugar.Warnw("[tts] edge-tts not found in PATH", "error", err)
		return false
	}
	return true
}

var DefaultVoices = struct {
	Uz string
	Ru string
}{
	Uz: "en-US-JennyNeural",
	Ru: "en-US-JennyNeural",
}

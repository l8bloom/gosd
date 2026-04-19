package gosd

import (
	"os"
	"testing"
)

func TestNewUpscalerCtx(t *testing.T) {
	Load()
	ctx := NewUpscalerCtx(
		os.Getenv("UPSCALER_PATH"),
		false,
		false,
		8,
		64,
	)

	if ctx == 0 {
		t.Error("Expected upscaler's context to be initialized, got nil pointer.")
		t.Log(ctx)
	}

	FreeUpscalerCtx(ctx)
}

func TestGetUpscaleFactor(t *testing.T) {
	ctx := NewUpscalerCtx(
		os.Getenv("UPSCALER_PATH"),
		false,
		false,
		8,
		64,
	)

	if ctx == 0 {
		t.Error("Expected upscaler's context to be initialized, got nil pointer.")
		t.Log(ctx)
	}

	factor := GetUpscaleFactor(ctx)
	// Check for a realistic multiplier, it is model dependent
	if factor <= 1 || factor > 16 {
		t.Errorf("received unrealistic upscale factor: %d. Is the model path correct?", factor)
	}

	FreeUpscalerCtx(ctx)
}

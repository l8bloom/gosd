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

func TestUpscale(t *testing.T) {
	imgParams := ImageGenParamsInit()
	imgParams.Width = 32
	imgParams.Height = 32
	imgParams.SampleParams.SampleSteps = 1
	imgParams.Prompt = "An orange cat."

	ctxParams := ContextParamsInit()
	ctxParams.DiffusionModelPath = os.Getenv("DIFFUSION_MODEL_PATH")
	ctxParams.VAEPath = os.Getenv("VAE_PATH")
	ctxParams.LLMPath = os.Getenv("LLM_PATH")

	ctx := NewContext(ctxParams)
	if ctx == 0 {
		t.Error("Expected context to be initialized, got nil pointer.")
		t.Log(ctx)
	}

	image := GenerateImage(ctx, imgParams)
	FreeCtx(ctx)

	uctx := NewUpscalerCtx(
		os.Getenv("UPSCALER_PATH"),
		false,
		false,
		8,
		4,
	)

	if uctx == 0 {
		t.Error("Expected upscaler's context to be initialized, got nil pointer.")
		t.Log(ctx)
	}
	defer FreeUpscalerCtx(uctx)

	Upscale(uctx, image, 4)
}

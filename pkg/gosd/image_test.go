package gosd

import (
	"errors"
	"os"
	"testing"
)

// test only some sensible default values
func TestImageGenParamsInit(t *testing.T) {
	Load()
	imgParams := ImageGenParamsInit()

	if imgParams.Width != 512 {
		t.Errorf("expected Width=512, got %d", imgParams.Width)
	}
	if imgParams.Height != 512 {
		t.Errorf("expected Height=512, got %d", imgParams.Height)
	}
	if imgParams.VAETilingParams.Enabled {
		t.Error("expected VAE Tiling disabled, got enabled")
	}
	if imgParams.VAETilingParams.TileSizeX != 0 {
		t.Errorf("expected VAE TileSizeX=0, got %d", imgParams.VAETilingParams.TileSizeX)
	}
	if imgParams.VAETilingParams.TileSizeY != 0 {
		t.Errorf("expected VAE TileSizeY=0, got %d", imgParams.VAETilingParams.TileSizeY)
	}
	if imgParams.VAETilingParams.TargetOverlap != 0.5 {
		t.Errorf("expected VAE TargetOverlap=0.5, got %f", imgParams.VAETilingParams.TargetOverlap)
	}
	if imgParams.VAETilingParams.RelSizeX != 0 {
		t.Errorf("expected VAE RelSizeX=0, got %f", imgParams.VAETilingParams.RelSizeX)
	}
	if imgParams.VAETilingParams.RelSizeY != 0 {
		t.Errorf("expected VAE RelSizeY=0, got %f", imgParams.VAETilingParams.RelSizeY)
	}
}

func TestImageGenParamsToStr(t *testing.T) {
	imgParams := ImageGenParamsInit()
	imgParamsStr := ImageGenParamsToStr(imgParams)
	if len(imgParamsStr) == 0 {
		t.Errorf("Expected non-empty image params string, got %s", imgParamsStr)
	}
}

func TestGenerateImage(t *testing.T) {
	imgParams := ImageGenParamsInit()
	imgParams.Width = 64
	imgParams.Height = 64
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

	defer FreeCtx(ctx)

	image := GenerateImage(ctx, imgParams)
	image.SavePNG("test_output.png")
	_, err := os.Stat("test_output.png")
	if errors.Is(err, os.ErrNotExist) {
		t.Error("the generated test image has not been saved")
	}
	if image.Width != 64 {
		t.Errorf("Expected image width=64, got %d", image.Width)
	}
	if image.Height != 64 {
		t.Errorf("Expected image height=64, got %d", image.Height)
	}
	if image.Channel != 3 {
		t.Errorf("Expected image channels=3, got %d", image.Channel)
	}
	if len(image.Data) != 64*64*3 {
		t.Error("the image data content should be 64x64x3")
	}
	os.Remove("test_output.png")
}

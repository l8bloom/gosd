package gosd

import (
	"errors"
	"os"
	"testing"
)

// test only some sensible default values
func TestVideoGenParamsInit(t *testing.T) {
	Load()
	vidParams := VideoGenParamsInit()

	if vidParams.Width != 512 {
		t.Errorf("expected Width=512, got %d", vidParams.Width)
	}
	if vidParams.Height != 512 {
		t.Errorf("expected Height=512, got %d", vidParams.Height)
	}
	if vidParams.VAETilingParams.Enabled {
		t.Error("expected VAE Tiling disabled, got enabled")
	}
	if vidParams.VAETilingParams.TileSizeX != 0 {
		t.Errorf("expected VAE TileSizeX=0, got %d", vidParams.VAETilingParams.TileSizeX)
	}
	if vidParams.VAETilingParams.TileSizeY != 0 {
		t.Errorf("expected VAE TileSizeY=0, got %d", vidParams.VAETilingParams.TileSizeY)
	}
	if vidParams.VAETilingParams.TargetOverlap != 0.5 {
		t.Errorf("expected VAE TargetOverlap=0.5, got %f", vidParams.VAETilingParams.TargetOverlap)
	}
	if vidParams.VAETilingParams.RelSizeX != 0 {
		t.Errorf("expected VAE RelSizeX=0, got %f", vidParams.VAETilingParams.RelSizeX)
	}
	if vidParams.VAETilingParams.RelSizeY != 0 {
		t.Errorf("expected VAE RelSizeY=0, got %f", vidParams.VAETilingParams.RelSizeY)
	}
}

func TestGenerateVideo(t *testing.T) {
	Load()
	vidParams := VideoGenParamsInit()
	vidParams.Width = 64
	vidParams.Height = 64
	vidParams.Prompt = "An orange cat."
	vidParams.SampleParams.SampleSteps = 1
	vidParams.SampleParams.SampleMethod = EulerSampleMethod
	vidParams.SampleParams.Guidance.TextCfg = 6
	vidParams.HighNoiseSampleParams.SampleSteps = 1
	vidParams.VideoFrames = 1

	ctxParams := ContextParamsInit()
	ctxParams.DiffusionModelPath = os.Getenv("VIDEO_DIFFUSION_MODEL_PATH")
	ctxParams.VAEPath = os.Getenv("VIDEO_VAE_PATH")
	ctxParams.T5XXLPath = os.Getenv("VIDEO_T5XXL_PATH")

	ctx := NewContext(ctxParams)
	if ctx == 0 {
		t.Error("Expected context to be initialized, got nil pointer.")
		t.Log(ctx)
	}

	defer FreeCtx(ctx)

	video := GenerateVideo(ctx, vidParams)
	video.Save("test_output.mp4", 1)
	_, err := os.Stat("test_output.mp4")
	if errors.Is(err, os.ErrNotExist) {
		t.Error("the generated test video has not been saved")
	}
	if video.Data[0].Width != 64 {
		t.Errorf("Expected image width=64, got %d", video.Data[0].Width)
	}
	if video.Data[0].Height != 64 {
		t.Errorf("Expected image height=64, got %d", video.Data[0].Width)
	}
	if video.Data[0].Channel != 3 {
		t.Errorf("Expected image channels=3, got %d", video.Data[0].Width)
	}
	if len(video.Data) != int(vidParams.VideoFrames)*64*64*3 {
		t.Error("the image data content should be 64x64x3")
	}
	os.Remove("test_output.mp4")
}

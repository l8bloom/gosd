package gosd

import (
	"errors"
	"os"
	"testing"
	"unsafe"
)

// test only some sensible default values
func TestVideoGenParamsInit(t *testing.T) {
	if err := Load(); err != nil {
		t.Fatal(err.Error())
	}
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
	vidParams := VideoGenParamsInit()
	vidParams.Width = 32
	vidParams.Height = 32
	vidParams.Prompt = "An orange cat."
	vidParams.SampleParams.SampleMethod = EulerSampleMethod
	vidParams.SampleParams.Guidance.TextCfg = 6
	vidParams.SampleParams.SampleSteps = 1
	vidParams.VideoFrames = 10

	ctxParams := ContextParamsInit()
	ctxParams.DiffusionModelPath = os.Getenv("VIDEO_DIFFUSION_MODEL_PATH")
	ctxParams.VAEPath = os.Getenv("VIDEO_VAE_PATH")
	ctxParams.T5XXLPath = os.Getenv("VIDEO_T5XXL_PATH")
	ctxParams.DiffusionFlashAttn = true
	ctxParams.KeepClipOnCPU = true
	vidParams.VAETilingParams.Enabled = true
	vidParams.VAETilingParams.RelSizeX = 4
	vidParams.VAETilingParams.RelSizeY = 4

	ctx := NewContext(ctxParams)
	if ctx == 0 {
		t.Error("Expected context to be initialized, got nil pointer.")
		t.Log(ctx)
	}

	defer FreeCtx(ctx)

	// annul preview callback set by image tests
	SetPreviewCallback(
		func(step int32, image Video, isNoisy bool, data unsafe.Pointer) {
		},
		PreviewVAE,
		1,
		true,
		false,
		nil,
	)
	video := GenerateVideo(ctx, vidParams)
	video.Save("test_output.mp4", 2)
	_, err := os.Stat("test_output.mp4")
	if errors.Is(err, os.ErrNotExist) {
		t.Error("the generated test video has not been saved")
		t.Log(err)
	}
	if video.Data[0].Width != uint32(vidParams.Width) {
		t.Errorf("Expected image width=%d, got %d", vidParams.Width, video.Data[0].Width)
	}
	if video.Data[0].Height != uint32(vidParams.Height) {
		t.Errorf("Expected image height=%d, got %d", vidParams.Height, video.Data[0].Height)
	}
	if video.Data[0].Channel != 3 {
		t.Errorf("Expected image channels=3, got %d", video.Data[0].Channel)
	}
	// adding here +1 since it seems sd returns index of the last
	// frame element and not the total length
	// TODO: demistify it
	if len(video.Data)+1 != int(vidParams.VideoFrames) {
		t.Errorf("number of video frames should be %d, got %d", vidParams.VideoFrames, len(video.Data))
	}
	os.Remove("test_output.mp4")
}

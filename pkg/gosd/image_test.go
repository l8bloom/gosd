package gosd

import (
	"errors"
	"os"
	"testing"
	"unsafe"
)

// test only some sensible default values
func TestImageGenParamsInit(t *testing.T) {
	if err := Load(); err != nil {
		t.Fatal(err.Error())
	}
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

func myProgressCallback(step int32, steps int32, time float32, data unsafe.Pointer) {
	*(*int)(data)++
}

func myImagePreviewCallback(step int32, image Image, isNoisy bool, data unsafe.Pointer) {
	*(*int)(data)++
}

func myLogCallback(level LogLevel, text string, data unsafe.Pointer) {
	*(*int)(data)++
}

func TestGenerateImage(t *testing.T) {
	data := new(int)
	SetPreviewCallback(myImagePreviewCallback, PreviewVAE, 1, true, false, unsafe.Pointer(data))

	data2 := new(int)
	SetProgressCallback(myProgressCallback, unsafe.Pointer(data2))

	data3 := new(int)
	SetLogCallback(myLogCallback, unsafe.Pointer(data3))

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

	if !CtxSupportsImageGeneration(ctx) {
		t.Error("expected context to support image generation, but got False")
		t.Log(ctx)
	}

	image := GenerateImage(ctx, imgParams)
	if image.Width != uint32(imgParams.Width) {
		t.Errorf("Expected image width=%d, got %d", imgParams.Width, image.Width)
	}
	if image.Height != uint32(imgParams.Height) {
		t.Errorf("Expected image height=%d, got %d", imgParams.Height, image.Height)
	}
	if image.Channel != 3 {
		t.Errorf("Expected image channels=3, got %d", image.Channel)
	}
	if len(image.Data) != int(imgParams.Width*imgParams.Height)*3 {
		t.Errorf("the image data content should be %dx%dx3", imgParams.Width, imgParams.Height)
	}

	image.SavePNG("test_output.png")
	_, err := os.Stat("test_output.png")
	if errors.Is(err, os.ErrNotExist) {
		t.Error("the generated test image has not been saved")
	}
	os.Remove("test_output.png")

	if *data == 0 {
		t.Errorf("preview callback: data counter sentinel should be positive, got %d", *data)
	}

	if *data2 == 0 {
		t.Errorf("progress callback: data counter sentinel should be positive, got %d", *data2)
	}

	if *data3 == 0 {
		t.Errorf("log callback: data counter sentinel should be positive, got %d", *data3)
	}

	defaultMethod := SampleMethodName(GetDefaultSampleMethod(ctx))
	if defaultMethod != "euler" {
		t.Errorf("expected default sampler method to be `euler`, got %s", defaultMethod)
	}

	defaultScheduler := SchedulerName(GetDefaultScheduler(ctx, EulerSampleMethod))
	if defaultScheduler != "discrete" {
		t.Errorf("expected default scheduler method to be `discrete`, got %s", defaultScheduler)
	}

	FreeCtx(ctx)

	uctx := NewUpscalerCtx(
		os.Getenv("UPSCALER_PATH"),
		true,
		true,
		4,
		int(imgParams.Width),
	)

	if uctx == 0 {
		t.Error("expected upscaler's context to be initialized, got nil pointer.")
		t.Log(ctx)
	}
	defer FreeUpscalerCtx(uctx)

	Upscale(uctx, image, 4)
}

func TestHiresParamsInit(t *testing.T) {
	hiresParams := HiresParamsInit()

	if hiresParams.Enabled {
		t.Errorf("expected Enabled=false, got %t", hiresParams.Enabled)
	}
	if hiresParams.Upscaler != HiresUpscalerLatent {
		t.Errorf("expected Upscaler=HiresUpscalerLatent, got %d", hiresParams.Upscaler)
	}
	if hiresParams.ModelPath != "" {
		t.Errorf("expected ModelPath=\"\", got %s", hiresParams.ModelPath)
	}
	if hiresParams.Scale != 2.0 {
		t.Errorf("expected Scale=2.0, got %f", hiresParams.Scale)
	}
	if hiresParams.TargetHeight != 0 {
		t.Errorf("expected TargetHeight=0, got %d", hiresParams.TargetHeight)
	}
	if hiresParams.TargetWidth != 0 {
		t.Errorf("expected TargetWidth=0, got %d", hiresParams.TargetWidth)
	}
	if hiresParams.Steps != 0 {
		t.Errorf("expected Steps=0, got %d", hiresParams.Steps)
	}
	if hiresParams.DenoisingStrength != 0.7 {
		t.Errorf("expected DenoisingStrength=0.7, got %f", hiresParams.DenoisingStrength)
	}
	if hiresParams.UpscaleTileSize != 128 {
		t.Errorf("expected UpscaleTileSize=128, got %d", hiresParams.UpscaleTileSize)
	}

	hiresParams.toC() // call for now to please the golangci-lint
}

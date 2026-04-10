package gosd

import (
	"os"
	"testing"
	"unsafe"
)

func myLogCallback(level LogLevel, text string, data unsafe.Pointer) {
	*(*int)(data)++
}

func TestSetLogCallback(t *testing.T) {
	if err := Load(); err != nil {
		t.Fatal(err.Error())
	}
	data := new(int)
	SetLogCallback(myLogCallback, unsafe.Pointer(data))

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
		t.Error("expected context to be initialized, got nil pointer.")
		t.Log(ctx)
	}

	defer FreeCtx(ctx)
	if *data == 0 {
		t.Errorf("data counter sentinel should be positive, got %d", *data)
	}
}

func myProgressCallback(step int32, steps int32, time float32, data unsafe.Pointer) {
	*(*int)(data)++
}

func TestSetProgressCallback(t *testing.T) {
	data := new(int)
	SetProgressCallback(myProgressCallback, unsafe.Pointer(data))

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

	defer FreeCtx(ctx)
	if *data == 0 {
		t.Errorf("data counter sentinel should be positive, got %d", *data)
	}
}

func myImagePreviewCallback(step int32, image Image, isNoisy bool, data unsafe.Pointer) {
	*(*int)(data)++
}

func TestSetPreviewCallback(t *testing.T) {
	data := new(int)
	SetPreviewCallback(myImagePreviewCallback, PreviewVAE, 1, true, false, unsafe.Pointer(data))
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

	defer FreeCtx(ctx)

	GenerateImage(ctx, imgParams)
	if *data == 0 {
		t.Errorf("data counter sentinel should be positive, got %d", *data)
	}
}

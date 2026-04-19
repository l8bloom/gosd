package gosd

import (
	"os"
	"testing"
)

func TestSampleParamsInit(t *testing.T) {
	Load()
	sp := SampleParamsInit()

	if sp.SampleSteps != 20 {
		t.Errorf("expected SampleSteps=20, got %d", sp.SampleSteps)
	}

	if sp.Guidance.TextCfg != 7 {
		t.Errorf("expected Guidance.TextCfg=7, got %g", sp.Guidance.TextCfg)
	}

	if sp.Guidance.DistilledGuidance != 3.5 {
		t.Errorf("expected Guidance.DistilledGuidance=3.5, got %g", sp.Guidance.DistilledGuidance)
	}
}

func TestSampleParamsToStr(t *testing.T) {
	sp := SampleParamsInit()
	spStr := SampleParamsToStr(sp)

	if len(spStr) == 0 {
		t.Errorf("expected non-empty sampler params string, got %s", spStr)
	}
}

func TestGetDefaultSampleMethod(t *testing.T) {
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

	defaultMethod := SampleMethodName(GetDefaultSampleMethod(ctx))
	if defaultMethod != "euler" {
		t.Errorf("expected default sampler method to be `euler`, got %s", defaultMethod)
	}
}

func TestGetDefaultScheduler(t *testing.T) {
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

	defaultScheduler := SchedulerName(GetDefaultScheduler(ctx, EulerSampleMethod))
	if defaultScheduler != "discrete" {
		t.Errorf("expected default scheduler method to be `discrete`, got %s", defaultScheduler)
	}
}

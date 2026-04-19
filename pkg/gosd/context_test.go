package gosd

import (
	"testing"
)

// test only some sensible default values
func TestContextParamsInit(t *testing.T) {
	if err := Load(); err != nil {
		t.Fatal(err.Error())
	}
	cp := ContextParamsInit()
	if !cp.VAEDecodeOnly {
		t.Errorf("expected VAEDecodeOnly=true, got %t", cp.VAEDecodeOnly)
	}
	if !cp.FreeParamsImmediately {
		t.Errorf("expected FreeParamsImmediately=true, got %t", cp.FreeParamsImmediately)
	}
	if cp.RNG != CUDARNG {
		t.Errorf("expected cp.RNG=CUDARNG, got %d", cp.RNG)
	}
	if cp.LoraApplyMode != LoraApplyAuto {
		t.Errorf("expected LoraApplyMode=LoraApplyAuto, got %d", cp.LoraApplyMode)
	}
	if cp.ChromaT5MaskPad != 1 {
		t.Errorf("expected ChromaT5MaskPad=1, got %d", cp.ChromaT5MaskPad)
	}
}

func TestCtxParamsToStr(t *testing.T) {
	ctxParams := ContextParamsInit()
	ctxParamsStr := CtxParamsToStr(ctxParams)
	if len(ctxParamsStr) == 0 {
		t.Error("expected non-empty string representation of context params.")
		t.Log(ctxParamsStr)
	}
}

package gosd

import (
	"testing"
)

// test only some sensible default values
func TestContextParamsInit(t *testing.T) {
	Load()
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

// func TestNewContext(t *testing.T) {
// 	var context Context

// 	_ctxParams := ctxParams.toC()
// 	newContext.Call(unsafe.Pointer(&context), unsafe.Pointer(&_ctxParams))

// 	return context
// }

// func TestFreeCtx(t *testing.T) {
// 	freeCtx.Call(nil, unsafe.Pointer(&ctx))
// }

// func TestCtxParamsToStr(t *testing.T) {
// 	str := utilsGetNulString()

// 	_params := ctxParams.toC()
// 	ctxParamsToStr.Call(unsafe.Pointer(&str), unsafe.Pointer(&_params))

// 	return charToString(str)
// }

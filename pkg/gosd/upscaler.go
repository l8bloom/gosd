package gosd

import (
	"unsafe"

	"github.com/jupiterrider/ffi"
)

// opaque pointers
type (
	UpscalerContext uintptr
)

var (
	// SD_API upscaler_ctx_t* new_upscaler_ctx(const char* esrgan_path, bool offload_params_to_cpu, bool direct, int n_threads, int tile_size);
	newUpscalerCtx ffi.Fun

	// SD_API void free_upscaler_ctx(upscaler_ctx_t* upscaler_ctx);
	freeUpscalerCtx ffi.Fun

	// SD_API int get_upscale_factor(upscaler_ctx_t* upscaler_ctx);
	getUpscaleFactor ffi.Fun

	// SD_API sd_image_t upscale(upscaler_ctx_t* upscaler_ctx, sd_image_t input_image, uint32_t upscale_factor);
	upscale ffi.Fun

	// ffiTypeImage represents the C struct sd_image_t,
	// stable-diffusion.cpp upscale(and others) API expect value and not a pointer
	ffiTypeImage = ffi.NewType(
		&ffi.TypeUint32,
		&ffi.TypeUint32,
		&ffi.TypeUint32,
		&ffi.TypePointer,
	)
)

func loadUpscalerRoutines(lib ffi.Lib) error {
	var err error

	if newUpscalerCtx, err = lib.Prep(
		"new_upscaler_ctx",
		&ffi.TypePointer,
		&ffi.TypePointer,
		&ffi.TypeUint8,
		&ffi.TypeUint8,
		&ffi.TypeSint32,
		&ffi.TypeSint32,
	); err != nil {
		return loadError("new_upscaler_ctx", err)
	}

	if freeUpscalerCtx, err = lib.Prep(
		"free_upscaler_ctx",
		&ffi.TypeVoid,
		&ffi.TypePointer,
	); err != nil {
		return loadError("free_upscaler_ctx", err)
	}

	if getUpscaleFactor, err = lib.Prep(
		"get_upscale_factor",
		&ffi.TypeSint32,
		&ffi.TypePointer,
	); err != nil {
		return loadError("get_upscale_factor", err)
	}

	if upscale, err = lib.Prep(
		"upscale",
		&ffiTypeImage,
		&ffi.TypePointer,
		&ffiTypeImage,
		&ffi.TypeUint32,
	); err != nil {
		return loadError("upscale", err)
	}

	return nil
}

// NewUpscalerCtx creates context for the upscaler.
func NewUpscalerCtx(esrganPath string, offloadParamsToCPU bool, direct bool, nThreads int, tileSize int) UpscalerContext {
	var ctx UpscalerContext

	ep := stringToChar(esrganPath)
	offload := boolToByte(offloadParamsToCPU)
	_direct := boolToByte(direct)
	threadsCnt := int32(nThreads)
	tile := int32(tileSize)

	newUpscalerCtx.Call(
		unsafe.Pointer(&ctx),
		unsafe.Pointer(&ep),
		unsafe.Pointer(&offload),
		unsafe.Pointer(&_direct),
		unsafe.Pointer(&threadsCnt),
		unsafe.Pointer(&tile),
	)

	return ctx
}

// FreeUpscalerCtx deallocates memory reserved by the upscaler context.
func FreeUpscalerCtx(ctx UpscalerContext) {
	freeUpscalerCtx.Call(nil, unsafe.Pointer(&ctx))
}

// GetUpscaleFactor returns upscaler's factor.
func GetUpscaleFactor(ctx UpscalerContext) int {
	var res int32

	getUpscaleFactor.Call(unsafe.Pointer(&res), unsafe.Pointer(&ctx))

	return int(res)
}

// Upscale upscales provided image with ESRGAN.
func Upscale(ctx UpscalerContext, img Image, upscaleFactor uint) Image {
	var resImage image
	inImage := *img.toC()
	factor := uint32(upscaleFactor)

	upscale.Call(
		unsafe.Pointer(&resImage),
		unsafe.Pointer(&ctx),
		unsafe.Pointer(&inImage),
		unsafe.Pointer(&factor),
	)

	return *resImage.toGo()
}

package gosd

import (
	"unsafe"

	"github.com/jupiterrider/ffi"
)

var (
	// SD_API void sd_set_log_callback(sd_log_cb_t sd_log_cb, void* data);
	setLogCallback ffi.Fun

	// SD_API void sd_set_progress_callback(sd_progress_cb_t cb, void* data);
	setProgressCallback ffi.Fun

	// SD_API void sd_set_preview_callback(sd_preview_cb_t cb, enum preview_t mode, int interval, bool denoised, bool noisy, void* data);
	setPreviewCallback ffi.Fun
)

func loadCallbacks(lib ffi.Lib) error {
	var err error
	if setLogCallback, err = lib.Prep(
		"sd_set_log_callback",
		&ffi.TypeVoid,
		&ffi.TypePointer,
		&ffi.TypePointer,
	); err != nil {
		return loadError("sd_set_log_callback", err)
	}

	if setProgressCallback, err = lib.Prep(
		"sd_set_progress_callback",
		&ffi.TypeVoid,
		&ffi.TypePointer,
		&ffi.TypePointer,
	); err != nil {
		return loadError("sd_set_progress_callback", err)
	}

	if setPreviewCallback, err = lib.Prep(
		"sd_set_preview_callback",
		&ffi.TypeVoid,
		&ffi.TypePointer,
		&ffi.TypeSint32,
		&ffi.TypeSint32,
		&ffi.TypeUint8,
		&ffi.TypeUint8,
		&ffi.TypePointer,
	); err != nil {
		return loadError("sd_set_preview_callback", err)
	}
	return nil
}

// Type used for SetPreviewCallback representing generated data
type PreviewFrames interface {
	Image | Video
}

type LogLevel int32

const (
	Debug LogLevel = iota
	Info
	Warn
	Error
)

type PreviewMode int32

const (
	PreviewNone PreviewMode = iota
	PreviewPROJ
	PreviewTAE
	PreviewVAE
	PreviewCount
)

type LogCallback func(level LogLevel, text string, data unsafe.Pointer)

var logCallback unsafe.Pointer
var sizeOfClosure = unsafe.Sizeof(ffi.Closure{})

func SetLogCallback(callback LogCallback, data unsafe.Pointer) {
	if callback == nil {
		panic("Can't set nil as a callback")
	}

	closure := ffi.ClosureAlloc(sizeOfClosure, &logCallback)

	fn := ffi.NewCallback(func(cif *ffi.Cif, ret unsafe.Pointer, args *unsafe.Pointer, userData unsafe.Pointer) uintptr {
		if args == nil {
			return 1 // error
		}

		arg := unsafe.Slice(args, cif.NArgs)

		level := *(*LogLevel)(arg[0])
		text := *(**byte)(arg[1])
		data := *(*unsafe.Pointer)(arg[2])

		callback(
			level,
			charToString(text),
			data,
		)
		return 0
	})

	var cifCallback ffi.Cif
	if status := ffi.PrepCif(
		&cifCallback,
		ffi.DefaultAbi,
		3,
		&ffi.TypeVoid,
		&ffi.TypeSint32,
		&ffi.TypePointer,
		&ffi.TypePointer,
	); status != ffi.OK {
		panic(status)
	}

	if closure != nil {
		if status := ffi.PrepClosureLoc(closure, &cifCallback, fn, nil, logCallback); status != ffi.OK {
			panic(status)
		}
	}

	setLogCallback.Call(nil, &logCallback, unsafe.Pointer(&data))

}

type ProgressCallback func(step int32, steps int32, time float32, data unsafe.Pointer)

var progressCallback unsafe.Pointer

func SetProgressCallback(callback ProgressCallback, data unsafe.Pointer) {
	if callback == nil {
		panic("Can't set nil as a callback")
	}

	closure := ffi.ClosureAlloc(sizeOfClosure, &progressCallback)

	fn := ffi.NewCallback(func(cif *ffi.Cif, ret unsafe.Pointer, args *unsafe.Pointer, userData unsafe.Pointer) uintptr {
		if args == nil {
			return 1 // error
		}

		arg := unsafe.Slice(args, cif.NArgs)

		step := *(*int32)(arg[0])
		steps := *(*int32)(arg[1])
		time := *(*float32)(arg[2])
		data := *(*unsafe.Pointer)(arg[3])

		callback(
			step,
			steps,
			time,
			data,
		)
		return 0
	})

	var cifCallback ffi.Cif
	if status := ffi.PrepCif(&cifCallback,
		ffi.DefaultAbi,
		4,
		&ffi.TypeVoid,
		&ffi.TypeSint32,
		&ffi.TypeSint32,
		&ffi.TypeFloat,
		&ffi.TypePointer,
	); status != ffi.OK {
		panic(status)
	}

	if closure != nil {
		if status := ffi.PrepClosureLoc(closure, &cifCallback, fn, nil, progressCallback); status != ffi.OK {
			panic(status)
		}
	}

	setProgressCallback.Call(nil, &progressCallback, unsafe.Pointer(&data))

}

type PreviewCallback[T PreviewFrames] func(step int32, frames T, isNoisy bool, data unsafe.Pointer)

var previewCallback unsafe.Pointer // keep in global due to GC

// step: current iteration step
// previewMode: mode in which to do the preview
// interval: iteration step slider
// image: generated image passed from the stable diffusion
// denoised: Should preview denoised images
// noisy: Should preview noisy images
// data: any app data
func SetPreviewCallback[T PreviewFrames](callback PreviewCallback[T], previewMode PreviewMode, interval int32, denoised bool, noisy bool, data unsafe.Pointer) {
	if callback == nil {
		panic("Can't set nil as a callback")
	}

	closure := ffi.ClosureAlloc(sizeOfClosure, &previewCallback)

	fn := ffi.NewCallback(func(cif *ffi.Cif, ret unsafe.Pointer, args *unsafe.Pointer, userData unsafe.Pointer) uintptr {
		if args == nil {
			return 1 // error
		}

		arg := unsafe.Slice(args, cif.NArgs)

		step := *(*int32)(arg[0])
		frameCount := *(*int32)(arg[1])
		frames := *(**image)(arg[2])
		isNoisy := *(*bool)(arg[3])
		data := *(*unsafe.Pointer)(arg[4])

		if frameCount <= 0 {
			panic("SetPreviewCallback: Preview image has no frames to preview")
		}

		// split between video and image preview
		var previewFrames T

		switch any(previewFrames).(type) {
		case Image:
			if frameCount > 1 {
				panic("SetPreviewCallback: can't use Image to preview a video.(use type Video)")
			}
			_img := *frames.toGo()
			previewFrames = any(_img).(T)

		default:
			if frameCount == 1 {
				panic("SetPreviewCallback: can't use Video to preview an image.(use type Image)")
			}
			_Cvid := unsafe.Slice(frames, frameCount)
			_vid := make([]Image, 0, frameCount)
			for _, fr := range _Cvid {
				_vid = append(_vid, *fr.toGo())
			}
			_vidFrames := Video{Data: _vid}
			previewFrames = any(_vidFrames).(T)
		}

		callback(
			step,
			previewFrames,
			isNoisy,
			data,
		)
		return 0
	})
	var cifCallback ffi.Cif

	if status := ffi.PrepCif(&cifCallback,
		ffi.DefaultAbi,
		5,
		&ffi.TypeVoid,
		&ffi.TypeSint32,
		&ffi.TypeSint32,
		&ffi.TypePointer,
		&ffi.TypeUint8,
		&ffi.TypePointer,
	); status != ffi.OK {
		panic(status)
	}

	if closure != nil {
		if status := ffi.PrepClosureLoc(closure, &cifCallback, fn, nil, previewCallback); status != ffi.OK {
			panic(status)
		}
	}

	setPreviewCallback.Call(nil,
		&previewCallback,
		unsafe.Pointer(&previewMode),
		unsafe.Pointer(&interval),
		unsafe.Pointer(&denoised),
		unsafe.Pointer(&noisy),
		unsafe.Pointer(&data),
	)

}

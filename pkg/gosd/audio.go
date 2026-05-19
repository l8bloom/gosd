package gosd

import (
	"unsafe"

	"github.com/jupiterrider/ffi"
)

var (
	// SD_API void free_sd_audio(sd_audio_t* audio);
	freeAudio ffi.Fun
)

func loadAudioRoutines(lib ffi.Lib) error {
	var err error

	if freeAudio, err = lib.Prep(
		"free_sd_audio",
		&ffi.TypeVoid,
		&ffi.TypePointer,
	); err != nil {
		return loadError("free_sd_audio", err)
	}

	return nil
}

type audio struct {
	SampleRate  uint32   // uint32_t sample_rate
	Channels    uint32   // uint32_t channels
	SampleCount uint64   // uint64_t sample_count
	Data        *float32 // float *data
}

func (a *audio) toGo() *Audio {
	return &Audio{
		SampleRate:  a.SampleRate,
		Channels:    a.Channels,
		SampleCount: a.SampleCount,
		Data:        unsafe.Slice(a.Data, a.SampleCount*uint64(a.Channels)),

		// Ensapsulate C data here for later deallocation
		// gosd does not own this piece of memory, sd.cpp does
		audio: uintptr(unsafe.Pointer(a)),
	}
}

type Audio struct {
	SampleRate  uint32
	Channels    uint32
	SampleCount uint64
	Data        []float32

	audio uintptr
}

func (a *Audio) toC() *audio {
	var _data *float32

	if len(a.Data) > 0 {
		_data = &a.Data[0]
	}

	return &audio{
		SampleRate:  a.SampleRate,
		Channels:    a.Channels,
		SampleCount: a.SampleCount,
		Data:        _data,
	}
}

func (a *Audio) audioPtr() uintptr {
	return a.audio
}

func newAudio() *Audio {
	return &Audio{}
}

// FreeAudio frees up memory holding audio content of a generated video.
func FreeAudio(audio *Audio) {
	_audio := audio.audioPtr()

	freeAudio.Call(nil, unsafe.Pointer(&_audio))

	audio.SampleRate = 0
	audio.Channels = 0
	audio.SampleCount = 0
	audio.Data = nil
	audio.audio = 0
}

package gosd

import (
	"fmt"
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

// func (a *audio) toGo() *Audio {
// 	_data := make([]float32, a.SampleCount*uint64(a.Channels))
// 	copy(_data, unsafe.Slice(a.Data, a.SampleCount))

// 	return &Audio{
// 		SampleRate:  a.SampleRate,
// 		Channels:    a.Channels,
// 		SampleCount: a.SampleCount,
// 		Data:        _data,
// 	}
// }

type Audio struct {
	SampleRate  uint32
	Channels    uint32
	SampleCount uint64
	Data        []float32
}

func (a *Audio) toC() *audio {
	var _data *float32

	if len(a.Data) > 0 {
		fmt.Println("HERE")
		_data = &a.Data[0]
	}

	return &audio{
		SampleRate:  a.SampleRate,
		Channels:    a.Channels,
		SampleCount: a.SampleCount,
		Data:        _data,
	}
}

func newAudio() *Audio {
	return &Audio{}
}

// FreeAudio frees up memory holding audio content of a generated video.
func FreeAudio(audio Audio) {
	_audio := audio.toC()

	freeAudio.Call(nil, unsafe.Pointer(&_audio))
}

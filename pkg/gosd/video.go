package gosd

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"os/exec"
	"unsafe"

	"github.com/jupiterrider/ffi"
)

var (

	// SD_API void sd_vid_gen_params_init(sd_vid_gen_params_t* sd_vid_gen_params);
	vidGenParamsInit ffi.Fun

	// SD_API bool generate_video(sd_ctx_t* sd_ctx, const sd_vid_gen_params_t* sd_vid_gen_params, sd_image_t** frames_out, int* num_frames_out, sd_audio_t** audio_out);
	generateVideo ffi.Fun
)

func loadVideosRoutines(lib ffi.Lib) error {
	var err error
	if vidGenParamsInit, err = lib.Prep(
		"sd_vid_gen_params_init",
		&ffi.TypeVoid,
		&ffi.TypePointer,
	); err != nil {
		return loadError("sd_vid_gen_params_init", err)
	}

	if generateVideo, err = lib.Prep(
		"generate_video",
		&ffi.TypeUint8,
		&ffi.TypePointer,
		&ffi.TypePointer,
		&ffi.TypePointer,
		&ffi.TypePointer,
		&ffi.TypePointer,
	); err != nil {
		return loadError("generate_video", err)
	}

	return nil
}

type videoParams struct {
	Lora                  *loraType        // const sd_lora_t* loras;
	LoraCount             uint32           // uint32_t lora_count;
	Prompt                *byte            // const char* prompt;
	NegativePrompt        *byte            // const char* negative_prompt;
	ClipSkip              int32            // int clip_skip;
	InitImage             image            // sd_image_t init_image;
	EndImage              image            // sd_image_t end_image;
	ControlFrames         *image           // sd_image_t* control_frames;
	ControlFramesSize     int32            // int control_frames_size;
	Width                 int32            // int width;
	Height                int32            // int height;
	SampleParams          sampleParamsType // sd_sample_params_t sample_params;
	HighNoiseSampleParams sampleParamsType // sd_sample_params_t high_noise_sample_params;
	MOEBoundary           float32          // float moe_boundary;
	Strength              float32          // float strength;
	Seed                  int64            // int64_t seed;
	VideoFrames           int32            // int video_frames;
	FPS                   int32            // int fps;
	VACEStrength          float32          // float vace_strength;
	VAETilingParams       vAETilingParams  // sd_tiling_params_t vae_tiling_params;
	Cache                 cacheParams      // sd_cache_params_t cache;
}

func (vp *videoParams) toGo() *VideoParams {
	size := int(vp.LoraCount)
	newLora := make([]LoraType, 0, size)
	srcLora := unsafe.Slice(vp.Lora, size)
	for _, sl := range srcLora {
		newLora = append(newLora, *sl.toGo())
	}

	size = int(vp.ControlFramesSize)
	newFrames := make([]Image, 0, size)
	srcFrames := unsafe.Slice(vp.ControlFrames, size)
	for _, sf := range srcFrames {
		newFrames = append(newFrames, *sf.toGo())
	}
	return &VideoParams{
		Lora:                  newLora,
		LoraCount:             vp.LoraCount,
		Prompt:                charToString(vp.Prompt),
		NegativePrompt:        charToString(vp.NegativePrompt),
		ClipSkip:              vp.ClipSkip,
		InitImage:             *vp.InitImage.toGo(),
		EndImage:              *vp.EndImage.toGo(),
		ControlFrames:         newFrames,
		ControlFramesSize:     vp.ControlFramesSize,
		Width:                 vp.Width,
		Height:                vp.Height,
		SampleParams:          *vp.SampleParams.toGo(),
		HighNoiseSampleParams: *vp.HighNoiseSampleParams.toGo(),
		MOEBoundary:           vp.MOEBoundary,
		Strength:              vp.Strength,
		Seed:                  vp.Seed,
		VideoFrames:           vp.VideoFrames,
		FPS:                   vp.FPS,
		VACEStrength:          vp.VACEStrength,
		VAETilingParams:       *vp.VAETilingParams.toGo(),
		Cache:                 *vp.Cache.toGo(),
	}
}

type VideoParams struct {
	Lora                  []LoraType
	LoraCount             uint32
	Prompt                string
	NegativePrompt        string
	ClipSkip              int32
	InitImage             Image
	EndImage              Image
	ControlFrames         []Image
	ControlFramesSize     int32
	Width                 int32
	Height                int32
	SampleParams          SampleParamsType
	HighNoiseSampleParams SampleParamsType
	MOEBoundary           float32
	Strength              float32
	Seed                  int64
	VideoFrames           int32
	FPS                   int32
	VACEStrength          float32
	VAETilingParams       VAETilingParams
	Cache                 CacheParams
}

func (vp *VideoParams) toC() *videoParams {
	var _lora *loraType
	if vp.LoraCount != 0 {
		_lora = vp.Lora[0].toC()
	}

	var _ctlFrames *image
	if vp.ControlFramesSize != 0 {
		_ctlFrames = vp.ControlFrames[0].toC()
	}

	return &videoParams{
		Lora:                  _lora,
		LoraCount:             vp.LoraCount,
		Prompt:                stringToChar(vp.Prompt),
		NegativePrompt:        stringToChar(vp.NegativePrompt),
		ClipSkip:              vp.ClipSkip,
		InitImage:             *vp.InitImage.toC(),
		EndImage:              *vp.EndImage.toC(),
		ControlFrames:         _ctlFrames,
		ControlFramesSize:     vp.ControlFramesSize,
		Width:                 vp.Width,
		Height:                vp.Height,
		SampleParams:          *vp.SampleParams.toC(),
		HighNoiseSampleParams: *vp.HighNoiseSampleParams.toC(),
		MOEBoundary:           vp.MOEBoundary,
		Strength:              vp.Strength,
		Seed:                  vp.Seed,
		VideoFrames:           vp.VideoFrames,
		FPS:                   vp.FPS,
		VACEStrength:          vp.VACEStrength,
		VAETilingParams:       *vp.VAETilingParams.toC(),
		Cache:                 *vp.Cache.toC(),
	}
}

type Video struct {
	Data  []Image
	Audio Audio
}

// Save saves generated video to the local disk with the help of ffmpeg.
// NOTE: This is not a core feature of the library,
// just an example of what can be done with
// the generated video after stable diffusion finishes.
func (gv Video) Save(filename string, fps int) error {
	if len(gv.Audio.Data) == 0 {
		return gv.saveVideo(filename, fps)
	}

	return gv.saveVideoWithAudio(filename, fps)
}

// The implementation is sub-ideal but solid.
// It opens an additional pipe(simulating OS-level file descriptor)
// to stream audio signal to ffmpeg.
// Video frames are provided via OS stdin.
func (gv Video) saveVideoWithAudio(filename string, fps int) error {
	if len(gv.Data) == 0 {
		return fmt.Errorf("no video frames to save")
	}

	audioReader, audioWriter, err := os.Pipe()
	if err != nil {
		return err
	}

	defer func() {
		// just in case smth goes bad; otherwise ffmpeg hangs
		_ = audioReader.Close()
		_ = audioWriter.Close()
	}()

	cmd := exec.Command("ffmpeg",
		"-y", // Overwrite output

		// Video
		"-f", "rawvideo",
		"-vcodec", "rawvideo",
		"-pix_fmt", "rgba",
		"-s", fmt.Sprintf("%dx%d", gv.Data[0].Width, gv.Data[0].Height),
		"-r", fmt.Sprintf("%d", fps),
		"-i", "pipe:0", // Explicitly read video from stdin

		// Audio
		"-f", "f32le",
		"-ar", fmt.Sprintf("%d", gv.Audio.SampleRate),
		"-ac", fmt.Sprintf("%d", gv.Audio.Channels),
		"-i", "pipe:3", // Read audio from the 3rd extra file descriptor

		// Output Encoding Configuration
		"-c:v", "libx264",
		"-pix_fmt", "yuv420p",
		"-c:a", "aac",
		"-shortest", // End the video when the shortest input ends
		filename,
	)

	// Attach the audio pipe reader to the command as pipe:3
	cmd.ExtraFiles = []*os.File{audioReader}

	videoStdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	// Close the parent's copy of the pipe's reader(Go is not Rust)
	_ = audioReader.Close()

	errChan := make(chan error, 2)

	// Send video frames
	go func() {
		defer func() { _ = videoStdin.Close() }()
		for _, img := range gv.Data {
			if _, err := videoStdin.Write(img.pixelize().Pix); err != nil {
				errChan <- err
				return
			}
		}
		errChan <- nil
	}()

	// Send audio signals
	go func() {
		defer func() { _ = audioWriter.Close() }()

		buf := new(bytes.Buffer)
		err := binary.Write(buf, binary.LittleEndian, gv.Audio.Data)
		if err != nil {
			errChan <- err
			return
		}

		if _, err := audioWriter.Write(buf.Bytes()); err != nil {
			errChan <- err
			return
		}
		errChan <- nil
	}()

	// wait
	for range 2 {
		if err := <-errChan; err != nil {
			return err
		}
	}

	return cmd.Wait()
}

func (gv Video) saveVideo(filename string, fps int) error {
	cmd := exec.Command("ffmpeg",
		"-y", // Overwrite output

		"-f", "rawvideo", // Input format
		"-vcodec", "rawvideo",
		"-pix_fmt", "rgba", // assume RGBA
		"-s", fmt.Sprintf("%dx%d", gv.Data[0].Width, gv.Data[0].Height),
		"-r", fmt.Sprintf("%d", fps),

		"-i", "-", // Read from stdin
		"-c:v", "libx264", // H.264 codec
		"-pix_fmt", "yuv420p", // Standard pixel format for players
		filename,
	)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	for _, img := range gv.Data {
		if _, err := stdin.Write(img.pixelize().Pix); err != nil {
			return err
		}
	}

	if err := stdin.Close(); err != nil {
		return err
	}

	return cmd.Wait()
}

func newVideoParams() *videoParams {
	return &videoParams{
		Prompt:         utilsGetNulString(),
		NegativePrompt: utilsGetNulString(),
	}
}

// VideoGenParamsInit creates a set of default values for video generation.
func VideoGenParamsInit() VideoParams {
	params := newVideoParams()

	vidGenParamsInit.Call(nil, unsafe.Pointer(&params))

	return *params.toGo()
}

// GenerateVideo starts the inference loop for video generation.
func GenerateVideo(ctx Context, vidParams VideoParams) Video {
	var res uint8

	image := &image{}
	imgPtr := &image

	_vidParams := vidParams.toC()
	framesCnt := new(int32)

	audio := newAudio().toC()
	audioPtr := &audio

	generateVideo.Call(
		unsafe.Pointer(&res),
		unsafe.Pointer(&ctx),
		unsafe.Pointer(&_vidParams),
		unsafe.Pointer(&imgPtr),
		unsafe.Pointer(&framesCnt),
		unsafe.Pointer(&audioPtr),
	)

	if !byteToBool(res) {
		// panic for now
		panic("gosd: video generation failed")
	}

	gv := Video{}
	// Attach audio if present
	if *audioPtr != nil {
		gv.Audio = *(*audioPtr).toGo()
	}
	// Attach the frames
	images := unsafe.Slice(image, int(*framesCnt))
	gv.Data = make([]Image, 0, len(images))

	for _, img := range images {
		gv.Data = append(gv.Data, *img.toGo())
	}

	return gv
}

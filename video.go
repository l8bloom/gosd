package gosd

import (
	"fmt"
	"os/exec"
	"unsafe"

	"github.com/jupiterrider/ffi"
)

var (

	//SD_API void sd_vid_gen_params_init(sd_vid_gen_params_t* sd_vid_gen_params);
	vidGenParamsInit ffi.Fun

	//SD_API sd_image_t* generate_video(sd_ctx_t* sd_ctx, const sd_vid_gen_params_t* sd_vid_gen_params, int* num_frames_out);
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
		VACEStrength:          vp.VACEStrength,
		VAETilingParams:       *vp.VAETilingParams.toC(),
		Cache:                 *vp.Cache.toC(),
	}
}

type Video struct {
	Data []Image
}

// this is not a core feature of the library,
// just an example of what can be done with
// the generated video after stable diffusion finishes
func (gv Video) Save(filename string, fps int) error {
	// requires ffmpeg installed
	cmd := exec.Command("ffmpeg",
		"-y",             // Overwrite output
		"-f", "rawvideo", // Input format
		"-vcodec", "rawvideo",
		"-pix_fmt", "rgba", // Match your GeneratedImage format
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
		if _, err := stdin.Write(img.Pixelize().Pix); err != nil {
			return err
		}
	}

	defer func() {
		closeError := stdin.Close()
		if err == nil {
			err = closeError
		}
	}()

	return cmd.Wait()
}

func newVideoParams() *videoParams {
	return &videoParams{
		Prompt:         utilsGetNulString(),
		NegativePrompt: utilsGetNulString(),
	}
}

func VideoGenParamsInit() VideoParams {
	params := newVideoParams()

	vidGenParamsInit.Call(nil, unsafe.Pointer(&params))

	return *params.toGo()
}

func GenerateVideo(ctx Context, vidParams VideoParams) Video {
	image := &image{}
	_vidParams := vidParams.toC()
	framesCnt := new(int32)

	generateVideo.Call(
		unsafe.Pointer(&image),
		unsafe.Pointer(&ctx),
		unsafe.Pointer(&_vidParams),
		unsafe.Pointer(&framesCnt),
	)

	images := unsafe.Slice(image, int(*framesCnt))
	gv := Video{}
	for _, img := range images {
		gv.Data = append(gv.Data, *img.toGo())
	}

	return gv
}

package services

import (
	"image"
	"image/png"
	"os"
	"unsafe"

	"github.com/jupiterrider/ffi"
)

var (
	// SD_API sd_image_t* generate_image(sd_ctx_t* sd_ctx, const sd_img_gen_params_t* sd_img_gen_params);
	generateImage ffi.Fun

	// SD_API void sd_img_gen_params_init(sd_img_gen_params_t* sd_img_gen_params);
	imageGenParamsInit ffi.Fun
)

func loadImageRoutines(lib ffi.Lib) error {
	var err error
	if generateImage, err = lib.Prep("generate_image", &ffi.TypePointer, &ffi.TypePointer, &ffi.TypePointer); err != nil {
		return loadError("generate_image", err)
	}
	if imageGenParamsInit, err = lib.Prep("sd_img_gen_params_init", &ffi.TypeVoid, &ffi.TypePointer); err != nil {
		return loadError("sd_img_gen_params_init", err)
	}

	return nil
}

// sd_lora_t
type LoraType struct {
	IsHighNoise uint8   // bool is_high_noise;
	Multiplier  float32 // float multiplier;
	Path        *byte   // const char* path;
}

// avoid using this and use [GeneratedImage] instead.
// [Image] depends on C type system, and may cause
// unwanted Go effects when GC deployed
type Image struct {
	Width   uint32 // uint32_t width
	Height  uint32 // uint32_t height
	Channel uint32 // uint32_t channel
	Data    *uint8 // uint8_t  *data
}

type SchedulerType int32

const (
	DiscreteScheduler SchedulerType = iota
	KarrasScheduler
	ExponentialScheduler
	AYSScheduler
	GITSScheduler
	SGMUniformScheduler
	SimpleScheduler
	SmoothStepScheduler
	KLOptimalScheduler
	LCMScheduler
	BongTangentScheduler
	SchedulerCount
)

type SampleMethodType int32

const (
	EulerSampleMethod SampleMethodType = iota
	EulerASampleMethod
	HeunSampleMethod
	DPM2SampleMethod
	DPMPP2SASampleMethod
	DPMPP2MSampleMethod
	DPMPP2Mv2SampleMethod
	IPNDMSampleMethod
	IPNDMVSampleMethod
	LCMSampleMethod
	DDIMTrailingSampleMethod
	TCDSampleMethod
	RESMultistepSampleMethod
	RES2SSampleMethod
	SampleMethodCount
)

type GuidanceParams struct {
	TextCfg           float32   // float txt_cfg;
	ImageCfg          float32   // float img_cfg;
	DistilledGuidance float32   // float distilled_guidance;
	SLG               SLGParams // sd_slg_params_t slg;
}

type SLGParams struct {
	Layers     *int32  // int* layers;
	LayerCount uint64  // size_t layer_count;
	LayerStart float32 // float layer_start;
	LayerEnd   float32 // float layer_end;
	Scale      float32 // float scale;
}

type SampleParamsType struct {
	Guidance          GuidanceParams   // sd_guidance_params_t guidance;
	Scheduler         SchedulerType    // enum scheduler_t scheduler;
	SampleMethod      SampleMethodType // enum sample_method_t sample_method;
	SampleSteps       int32            // int sample_steps;
	ETA               float32          // float eta;
	ShiftedTimestamp  int32            // int shifted_timestep;
	CustomSigmas      *float32         // float* custom_sigmas;
	CustomSigmasCount int32            // int custom_sigmas_count;
	FlowShift         float32          // float flow_shift;
}

type PMParamsType struct {
	IDImages      *Image  // sd_image_t* id_images;
	IDImagesCount int32   // int id_images_count;
	IDEmbedPath   *byte   // const char* id_embed_path;
	StyleStrength float32 // float style_strength;
}

type VAETilingParams struct {
	Enabled       uint8   // bool enabled;
	TileSizeX     int32   // int tile_size_x;
	TileSizeY     int32   // int tile_size_y;
	TargetOverlap float32 // float target_overlap;
	RelSizeX      float32 // float rel_size_x;
	RelSizeY      float32 // float rel_size_y;
}

type CacheModeType int32

const (
	CacheDisabled CacheModeType = iota
	CacheEasyCache
	CacheUcache
	CacheDBcache
	CacheTaylorseer
	CacheCacheDit
)

type CacheParams struct {
	Mode                     CacheModeType // enum sd_cache_mode_t mode;
	ReuseThreshold           float32       // float reuse_threshold;
	StartPercent             float32       // float start_percent;
	EndPercent               float32       // float end_percent;
	ErrorDecayRate           float32       // float error_decay_rate;
	UseRelativeThreshold     uint8         // bool use_relative_threshold;
	ResetErrorOnCompute      uint8         // bool reset_error_on_compute;
	FNComputeBlocks          int32         // int Fn_compute_blocks;
	BNComputeBlocks          int32         // int Bn_compute_blocks;
	ResidualDiffThreshold    float32       // float residual_diff_threshold;
	MaxWarmupSteps           int32         // int max_warmup_steps;
	MaxCachedSteps           int32         // int max_cached_steps;
	MaxContinuousCachedSteps int32         // int max_continuous_cached_steps;
	TaylorSeerNDerivatives   int32         // int taylorseer_n_derivatives;
	TaylorSeerSkipInterval   int32         // int taylorseer_skip_interval;
	SCMMask                  *byte         // const char* scm_mask;
	SCMPolicyDynamic         uint8         // bool scm_policy_dynamic;
}

type ImageParams struct {
	Lora               *LoraType        // const sd_lora_t* loras;
	LoraCount          uint32           // uint32_t lora_count;
	Prompt             *byte            // const char* prompt;
	NegativePrompt     *byte            // const char* negative_prompt;
	ClipSkip           int32            // int clip_skip;
	InitImage          Image            // sd_image_t init_image;
	RefImages          *Image           // sd_image_t* ref_images;
	RefImagesCount     int32            // int ref_images_count;
	AutoResizeRefImage uint8            // bool auto_resize_ref_image;
	IncreaseRefIndex   uint8            // bool increase_ref_index;
	MaskImage          Image            // sd_image_t mask_image;
	Width              int32            // int width;
	Height             int32            // int height;
	SampleParams       SampleParamsType // sd_sample_params_t sample_params;
	Strength           float32          // float strength;
	Seed               int64            // int64_t seed;
	BatchCount         int32            // int batch_count;
	ControlImage       Image            // sd_image_t control_image;
	ControlStrength    float32          // float control_strength;
	PMParams           PMParamsType     // sd_pm_params_t pm_params;
	VAETilingParams    VAETilingParams  // sd_tiling_params_t vae_tiling_params;
	Cache              CacheParams      // sd_cache_params_t cache;
}

func GenerateImage(ctx Context, ip ImageParams) GeneratedImage {
	var image *Image

	i := &ip
	generateImage.Call(unsafe.Pointer(&image), unsafe.Pointer(&ctx), unsafe.Pointer(&i))

	return NewGeneratedImage(*image)
}

func NewImageParams() *ImageParams {
	ip := &ImageParams{
		Prompt:         utilsGetNulString(),
		NegativePrompt: utilsGetNulString(),
	}

	return ip
}

func ImageGenParamsInit() ImageParams {
	var ip *ImageParams = NewImageParams()

	imageGenParamsInit.Call(nil, unsafe.Pointer(&ip))
	return *ip
}

type GeneratedImage struct {
	Width   int
	Height  int
	Channel int
	Data    []uint8
}

func (img GeneratedImage) SavePNG(filename string) error {
	if len(img.Data) == 0 {
		panic("Image with 0 length.")
	}

	rect := image.Rect(0, 0, img.Width, img.Height)
	rgba := image.NewRGBA(rect)

	src := img.Data
	dst := rgba.Pix
	channels := img.Channel

	// Single loop that handles 1-channel (Grey) or 3-channel (RGB)
	for i, j := 0, 0; i < len(src); i += channels {
		if channels == 3 {
			dst[j] = src[i]     // R
			dst[j+1] = src[i+1] // G
			dst[j+2] = src[i+2] // B
		} else if channels == 1 {
			// Monochromatic: Map Gray to R, G, and B
			val := src[i]
			dst[j] = val
			dst[j+1] = val
			dst[j+2] = val
		}
		dst[j+3] = 255 // Alpha is ALWAYS required for RGBA
		j += 4
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, rgba)
}

func NewGeneratedImage(img Image) GeneratedImage {
	size := int(img.Width * img.Height * img.Channel)
	newPix := make([]uint8, size)

	// Cast to Go slice
	srcPix := unsafe.Slice(img.Data, size)
	copy(newPix, srcPix)

	return GeneratedImage{
		Width:   int(img.Width),
		Height:  int(img.Height),
		Channel: int(img.Channel),
		Data:    newPix,
	}
}

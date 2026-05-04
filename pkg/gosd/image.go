package gosd

import (
	imgPckg "image"
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

	// SD_API char* sd_img_gen_params_to_str(const sd_img_gen_params_t* sd_img_gen_params);
	imageGenParamsToStr ffi.Fun

	// SD_API void sd_hires_params_init(sd_hires_params_t* hires_params);
	hiresParamsInit ffi.Fun
)

func loadImageRoutines(lib ffi.Lib) error {
	var err error
	if generateImage, err = lib.Prep(
		"generate_image",
		&ffi.TypePointer,
		&ffi.TypePointer,
		&ffi.TypePointer,
	); err != nil {
		return loadError("generate_image", err)
	}
	if imageGenParamsInit, err = lib.Prep(
		"sd_img_gen_params_init",
		&ffi.TypeVoid,
		&ffi.TypePointer,
	); err != nil {
		return loadError("sd_img_gen_params_init", err)
	}

	if imageGenParamsToStr, err = lib.Prep(
		"sd_img_gen_params_to_str",
		&ffi.TypePointer,
		&ffi.TypePointer,
	); err != nil {
		return loadError("sd_img_gen_params_to_str", err)
	}

	if hiresParamsInit, err = lib.Prep(
		"sd_hires_params_init",
		&ffi.TypeVoid,
		&ffi.TypePointer,
	); err != nil {
		return loadError("sd_hires_params_init", err)
	}

	return nil
}

// sd_lora_t
type loraType struct {
	IsHighNoise uint8   // bool is_high_noise;
	Multiplier  float32 // float multiplier;
	Path        *byte   // const char* path;
}

func (lt *loraType) toGo() *LoraType {
	if lt == nil {
		return nil
	}
	return &LoraType{
		IsHighNoise: byteToBool(lt.IsHighNoise),
		Multiplier:  lt.Multiplier,
		Path:        charToString(lt.Path),
	}
}

type LoraType struct {
	IsHighNoise bool
	Multiplier  float32
	Path        string
}

func (lt *LoraType) toC() *loraType {
	return &loraType{
		IsHighNoise: boolToByte(lt.IsHighNoise),
		Multiplier:  lt.Multiplier,
		Path:        stringToChar(lt.Path),
	}
}

type image struct {
	Width   uint32 // uint32_t width
	Height  uint32 // uint32_t height
	Channel uint32 // uint32_t channel
	Data    *uint8 // uint8_t  *data
}

func (lt *image) toGo() *Image {
	size := int(lt.Width * lt.Height * lt.Channel)
	newPix := make([]uint8, size)

	// Cast to Go slice
	srcPix := unsafe.Slice(lt.Data, size)
	copy(newPix, srcPix)

	return &Image{
		Width:   lt.Width,
		Height:  lt.Height,
		Channel: lt.Channel,
		Data:    newPix,
	}
}

type Image struct {
	Width   uint32
	Height  uint32
	Channel uint32
	Data    []uint8
}

func (lt *Image) toC() *image {
	size := int(lt.Width * lt.Height * lt.Channel)
	var _data *uint8

	if size != 0 {
		_data = &lt.Data[0]
	}

	return &image{
		Width:   lt.Width,
		Height:  lt.Height,
		Channel: lt.Channel,
		Data:    _data,
	}
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
	ERSDESampleMethod
	SampleMethodCount
)

type HiresUpscalerType int32

const (
	HiresUpscalerNone HiresUpscalerType = iota
	HiresUpscalerLatent
	HiresUpscalerLatentNearest
	HiresUpscalerLatentNearestExact
	HiresUpscalerLatentAntialiased
	HiresUpscalerLatentBicubic
	HiresUpscalerLatentBicubicAntialiased
	HiresUpscalerLanczos
	HiresUpscalerNearest
	HiresUpscalerModel
	HiresUpscalerCount
)

type hiresParams struct {
	Enabled           uint8             // bool enabled;
	Upscaler          HiresUpscalerType // enum sd_hires_upscaler_t upscaler;
	ModelPath         *byte             // const char* model_path;
	Scale             float32           // float scale;
	TargetWidth       int32             // int target_width;
	TargetHeight      int32             // int target_height;
	Steps             int32             // int steps;
	DenoisingStrength float32           // float denoising_strength;
	UpscaleTileSize   int32             // int upscale_tile_size;
}

func (hp *hiresParams) toGo() *HiresParams {
	return &HiresParams{
		Enabled:           byteToBool(hp.Enabled),
		Upscaler:          hp.Upscaler,
		ModelPath:         charToString(hp.ModelPath),
		Scale:             hp.Scale,
		TargetWidth:       hp.TargetWidth,
		TargetHeight:      hp.TargetHeight,
		Steps:             hp.Steps,
		DenoisingStrength: hp.DenoisingStrength,
		UpscaleTileSize:   hp.UpscaleTileSize,
	}
}

func newHiresParams() *hiresParams {
	hp := &hiresParams{
		ModelPath: utilsGetNulString(),
	}

	return hp
}

type HiresParams struct {
	Enabled           bool
	Upscaler          HiresUpscalerType
	ModelPath         string
	Scale             float32
	TargetWidth       int32
	TargetHeight      int32
	Steps             int32
	DenoisingStrength float32
	UpscaleTileSize   int32
}

func (hp *HiresParams) toC() *hiresParams {
	return &hiresParams{
		Enabled:           boolToByte(hp.Enabled),
		Upscaler:          hp.Upscaler,
		ModelPath:         stringToChar(hp.ModelPath),
		Scale:             hp.Scale,
		TargetWidth:       hp.TargetWidth,
		TargetHeight:      hp.TargetHeight,
		Steps:             hp.Steps,
		DenoisingStrength: hp.DenoisingStrength,
		UpscaleTileSize:   hp.UpscaleTileSize,
	}
}

type guidanceParams struct {
	TextCfg           float32   // float txt_cfg;
	ImageCfg          float32   // float img_cfg;
	DistilledGuidance float32   // float distilled_guidance;
	SLG               sLGParams // sd_slg_params_t slg;
}

func (g *guidanceParams) toGo() *GuidanceParams {
	return &GuidanceParams{
		TextCfg:           g.TextCfg,
		ImageCfg:          g.ImageCfg,
		DistilledGuidance: g.DistilledGuidance,
		SLG:               *g.SLG.toGo(),
	}
}

type GuidanceParams struct {
	TextCfg           float32
	ImageCfg          float32
	DistilledGuidance float32
	SLG               SLGParams
}

func (g *GuidanceParams) toC() *guidanceParams {
	return &guidanceParams{
		TextCfg:           g.TextCfg,
		ImageCfg:          g.ImageCfg,
		DistilledGuidance: g.DistilledGuidance,
		SLG:               *g.SLG.toC(),
	}
}

type sLGParams struct {
	Layers     *int32  // int* layers;
	LayerCount uint64  // size_t layer_count;
	LayerStart float32 // float layer_start;
	LayerEnd   float32 // float layer_end;
	Scale      float32 // float scale;
}

func (slg *sLGParams) toGo() *SLGParams {
	size := int(slg.LayerCount)
	newLayers := make([]int32, size)

	// Cast to Go slice
	srcLayers := unsafe.Slice(slg.Layers, size)
	copy(newLayers, srcLayers)

	return &SLGParams{
		Layers:     newLayers,
		LayerCount: slg.LayerCount,
		LayerStart: slg.LayerStart,
		LayerEnd:   slg.LayerEnd,
		Scale:      slg.Scale,
	}
}

type SLGParams struct {
	Layers     []int32
	LayerCount uint64
	LayerStart float32
	LayerEnd   float32
	Scale      float32
}

func (slg *SLGParams) toC() *sLGParams {
	size := int(slg.LayerCount)
	var _data *int32

	if size != 0 {
		_data = &slg.Layers[0]
	}
	return &sLGParams{
		Layers:     _data,
		LayerCount: slg.LayerCount,
		LayerStart: slg.LayerStart,
		LayerEnd:   slg.LayerEnd,
		Scale:      slg.Scale,
	}
}

type pMParamsType struct {
	IDImages      *image  // sd_image_t* id_images;
	IDImagesCount int32   // int id_images_count;
	IDEmbedPath   *byte   // const char* id_embed_path;
	StyleStrength float32 // float style_strength;
}

func (pmp *pMParamsType) toGo() *PMParamsType {
	size := int(pmp.IDImagesCount)
	newImages := make([]Image, size)

	srcImage := unsafe.Slice(pmp.IDImages, size)
	for _, sl := range srcImage {
		newImages = append(newImages, *sl.toGo())
	}

	return &PMParamsType{
		IDImages:      newImages,
		IDImagesCount: pmp.IDImagesCount,
		IDEmbedPath:   charToString(pmp.IDEmbedPath),
		StyleStrength: pmp.StyleStrength,
	}
}

type PMParamsType struct {
	IDImages      []Image
	IDImagesCount int32
	IDEmbedPath   string
	StyleStrength float32
}

func (pmp *PMParamsType) toC() *pMParamsType {
	size := int(pmp.IDImagesCount)
	var _data *image

	if size != 0 {
		_data = pmp.IDImages[10].toC()
	}

	return &pMParamsType{
		IDImages:      _data,
		IDImagesCount: pmp.IDImagesCount,
		IDEmbedPath:   stringToChar(pmp.IDEmbedPath),
		StyleStrength: pmp.StyleStrength,
	}
}

type vAETilingParams struct {
	Enabled       uint8   // bool enabled;
	TileSizeX     int32   // int tile_size_x;
	TileSizeY     int32   // int tile_size_y;
	TargetOverlap float32 // float target_overlap;
	RelSizeX      float32 // float rel_size_x;
	RelSizeY      float32 // float rel_size_y;
}

func (vae *vAETilingParams) toGo() *VAETilingParams {
	return &VAETilingParams{
		Enabled:       byteToBool(vae.Enabled),
		TileSizeX:     vae.TileSizeX,
		TileSizeY:     vae.TileSizeY,
		TargetOverlap: vae.TargetOverlap,
		RelSizeX:      vae.RelSizeX,
		RelSizeY:      vae.RelSizeY,
	}
}

type VAETilingParams struct {
	Enabled       bool
	TileSizeX     int32
	TileSizeY     int32
	TargetOverlap float32
	RelSizeX      float32
	RelSizeY      float32
}

func (vae *VAETilingParams) toC() *vAETilingParams {
	return &vAETilingParams{
		Enabled:       boolToByte(vae.Enabled),
		TileSizeX:     vae.TileSizeX,
		TileSizeY:     vae.TileSizeY,
		TargetOverlap: vae.TargetOverlap,
		RelSizeX:      vae.RelSizeX,
		RelSizeY:      vae.RelSizeY,
	}
}

type imageParams struct {
	Lora               *loraType        // const sd_lora_t* loras;
	LoraCount          uint32           // uint32_t lora_count;
	Prompt             *byte            // const char* prompt;
	NegativePrompt     *byte            // const char* negative_prompt;
	ClipSkip           int32            // int clip_skip;
	InitImage          image            // sd_image_t init_image;
	RefImages          *image           // sd_image_t* ref_images;
	RefImagesCount     int32            // int ref_images_count;
	AutoResizeRefImage uint8            // bool auto_resize_ref_image;
	IncreaseRefIndex   uint8            // bool increase_ref_index;
	MaskImage          image            // sd_image_t mask_image;
	Width              int32            // int width;
	Height             int32            // int height;
	SampleParams       sampleParamsType // sd_sample_params_t sample_params;
	Strength           float32          // float strength;
	Seed               int64            // int64_t seed;
	BatchCount         int32            // int batch_count;
	ControlImage       image            // sd_image_t control_image;
	ControlStrength    float32          // float control_strength;
	PMParams           pMParamsType     // sd_pm_params_t pm_params;
	VAETilingParams    vAETilingParams  // sd_tiling_params_t vae_tiling_params;
	Cache              cacheParams      // sd_cache_params_t cache;
	HiresParams        hiresParams      // sd_hires_params_t hires;
}

func (i *imageParams) toGo() *ImageParams {
	size := int(i.LoraCount)
	newLora := make([]LoraType, 0, size)

	// Cast to Go slice
	srcLora := unsafe.Slice(i.Lora, size)
	for _, sl := range srcLora {
		newLora = append(newLora, *sl.toGo())
	}

	size = int(i.RefImagesCount)
	newImages := make([]Image, 0, size)

	// Cast to Go slice
	srcImages := unsafe.Slice(i.RefImages, size)
	for _, si := range srcImages {
		newImages = append(newImages, *si.toGo())
	}

	return &ImageParams{
		Lora:               newLora,
		LoraCount:          i.LoraCount,
		Prompt:             charToString(i.Prompt),
		NegativePrompt:     charToString(i.NegativePrompt),
		ClipSkip:           i.ClipSkip,
		InitImage:          *i.InitImage.toGo(),
		RefImages:          newImages,
		RefImagesCount:     i.RefImagesCount,
		AutoResizeRefImage: byteToBool(i.AutoResizeRefImage),
		IncreaseRefIndex:   byteToBool(i.IncreaseRefIndex),
		MaskImage:          *i.MaskImage.toGo(),
		Width:              i.Width,
		Height:             i.Height,
		SampleParams:       *i.SampleParams.toGo(),
		Strength:           i.Strength,
		Seed:               i.Seed,
		BatchCount:         i.BatchCount,
		ControlImage:       *i.ControlImage.toGo(),
		ControlStrength:    i.ControlStrength,
		PMParams:           *i.PMParams.toGo(),
		VAETilingParams:    *i.VAETilingParams.toGo(),
		Cache:              *i.Cache.toGo(),
		HiresParams:        *i.HiresParams.toGo(),
	}
}

type ImageParams struct {
	Lora               []LoraType
	LoraCount          uint32
	Prompt             string
	NegativePrompt     string
	ClipSkip           int32
	InitImage          Image
	RefImages          []Image
	RefImagesCount     int32
	AutoResizeRefImage bool
	IncreaseRefIndex   bool
	MaskImage          Image
	Width              int32
	Height             int32
	SampleParams       SampleParamsType
	Strength           float32
	Seed               int64
	BatchCount         int32
	ControlImage       Image
	ControlStrength    float32
	PMParams           PMParamsType
	VAETilingParams    VAETilingParams
	Cache              CacheParams
	HiresParams        HiresParams
}

func (i *ImageParams) toC() *imageParams {
	var _lora *loraType
	if i.LoraCount != 0 {
		_lora = i.Lora[0].toC()
	}
	var _refImages *image
	if i.RefImagesCount != 0 {
		_refImages = i.RefImages[0].toC()
	}

	return &imageParams{
		Lora:               _lora,
		LoraCount:          i.LoraCount,
		Prompt:             stringToChar(i.Prompt),
		NegativePrompt:     stringToChar(i.NegativePrompt),
		ClipSkip:           i.ClipSkip,
		InitImage:          *i.InitImage.toC(),
		RefImages:          _refImages,
		RefImagesCount:     i.RefImagesCount,
		AutoResizeRefImage: boolToByte(i.AutoResizeRefImage),
		IncreaseRefIndex:   boolToByte(i.IncreaseRefIndex),
		MaskImage:          *i.MaskImage.toC(),
		Width:              i.Width,
		Height:             i.Height,
		SampleParams:       *i.SampleParams.toC(),
		Strength:           i.Strength,
		Seed:               i.Seed,
		BatchCount:         i.BatchCount,
		ControlImage:       *i.ControlImage.toC(),
		ControlStrength:    i.ControlStrength,
		PMParams:           *i.PMParams.toC(),
		VAETilingParams:    *i.VAETilingParams.toC(),
		Cache:              *i.Cache.toC(),
		HiresParams:        *i.HiresParams.toC(),
	}
}

func GenerateImage(ctx Context, ip ImageParams) Image {
	var image *image

	i := ip.toC()
	generateImage.Call(unsafe.Pointer(&image), unsafe.Pointer(&ctx), unsafe.Pointer(&i))

	return *image.toGo()
}

func newImageParams() *imageParams {
	ip := &imageParams{
		Prompt:         utilsGetNulString(),
		NegativePrompt: utilsGetNulString(),
	}

	return ip
}

func ImageGenParamsInit() ImageParams {
	ip := newImageParams()

	imageGenParamsInit.Call(nil, unsafe.Pointer(&ip))
	return *ip.toGo()
}

func ImageGenParamsToStr(ip ImageParams) string {
	str := utilsGetNulString()

	_params := ip.toC()
	imageGenParamsToStr.Call(unsafe.Pointer(&str), unsafe.Pointer(&_params))

	return charToString(str)
}

// this is not a core feature of the library,
// just an example of what can be done with
// the generated image from the stable diffusion
func (img Image) SavePNG(filename string) error {
	pix := img.Pixelize()

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func() {
		closeError := f.Close()
		if err == nil {
			err = closeError
		}
	}()

	return png.Encode(f, &pix)
}

func (img Image) Pixelize() imgPckg.RGBA {
	if len(img.Data) == 0 {
		panic("Image with 0 length.")
	}

	rect := imgPckg.Rect(0, 0, int(img.Width), int(img.Height))
	rgba := imgPckg.NewRGBA(rect)

	src := img.Data
	dst := rgba.Pix
	channels := int(img.Channel)

	// Tries to handle gray and colorful images
	for i, j := 0, 0; i < len(src); i += channels {
		switch channels {
		case 3:
			dst[j] = src[i]     // R
			dst[j+1] = src[i+1] // G
			dst[j+2] = src[i+2] // B
		default:
			// Monochromatic
			val := src[i]
			dst[j] = val
			dst[j+1] = val
			dst[j+2] = val
		}
		dst[j+3] = 255 // Alpha
		j += 4
	}
	return *rgba
}

func HiresParamsInit() HiresParams {
	hp := newHiresParams()

	hiresParamsInit.Call(nil, unsafe.Pointer(&hp))
	return *hp.toGo()
}

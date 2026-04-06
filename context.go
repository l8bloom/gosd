package gosd

import (
	"unsafe"

	"github.com/jupiterrider/ffi"
)

var (
	// SD_API void sd_ctx_params_init(sd_ctx_params_t* sd_ctx_params);
	contextParamsInit ffi.Fun

	// SD_API sd_ctx_t* new_sd_ctx(const sd_ctx_params_t* sd_ctx_params);
	newContext ffi.Fun

	// SD_API char* sd_ctx_params_to_str(const sd_ctx_params_t* sd_ctx_params);
	ctxParamsToStr ffi.Fun

	// SD_API void free_sd_ctx(sd_ctx_t* sd_ctx);
	freeCtx ffi.Fun
)

func loadContextRoutines(lib ffi.Lib) error {
	var err error
	if newContext, err = lib.Prep(
		"new_sd_ctx",
		&ffi.TypePointer,
		&ffi.TypePointer,
	); err != nil {
		return loadError("new_sd_ctx", err)
	}

	if contextParamsInit, err = lib.Prep(
		"sd_ctx_params_init",
		&ffi.TypeVoid,
		&ffi.TypePointer,
	); err != nil {
		return loadError("sd_ctx_params_init", err)
	}

	if freeCtx, err = lib.Prep(
		"free_sd_ctx",
		&ffi.TypeVoid,
		&ffi.TypePointer,
	); err != nil {
		return loadError("free_sd_ctx", err)
	}

	if ctxParamsToStr, err = lib.Prep(
		"sd_ctx_params_to_str",
		&ffi.TypePointer,
		&ffi.TypePointer,
	); err != nil {
		return loadError("sd_ctx_params_to_str", err)
	}

	return nil
}

// opaque pointers
type (
	Context uintptr
)

// C type
type contextParams struct {
	ModelPath                   *byte             // const char* model_path;
	ClipLPath                   *byte             // const char* clip_l_path;
	ClipGPath                   *byte             // const char* clip_g_path;
	ClipVisionPath              *byte             // const char* clip_vision_path;
	T5XXLPath                   *byte             // const char* t5xxl_path;
	LLMPath                     *byte             // const char* llm_path;
	LLMVisionPath               *byte             // const char* llm_vision_path;
	DiffusionModelPath          *byte             // const char* diffusion_model_path;
	HighNoiseDiffusionModelPath *byte             // const char* high_noise_diffusion_model_path;
	VAEPath                     *byte             // const char* vae_path;
	TAESDPath                   *byte             // const char* taesd_path;
	ControlNetPath              *byte             // const char* control_net_path;
	Embeddings                  *embedding        // const sd_embedding_t* embeddings;
	EmbeddingCount              uint32            // uint32_t embedding_count;
	PhotoMakerPath              *byte             // const char* photo_maker_path;
	TensorTypeRules             *byte             // const char* tensor_type_rules;
	VAEDecodeOnly               uint8             // bool vae_decode_only;
	FreeParamsImmediately       uint8             // bool free_params_immediately;
	NThreads                    int32             // int n_threads;
	WType                       SDType            // enum sd_type_t wtype;
	RNG                         RNGType           // enum rng_type_t rng_type;
	SamplerRNG                  RNGType           // enum rng_type_t sampler_rng_type;
	Prediction                  PredictionType    // enum prediction_t prediction;
	LoraApplyMode               LoraApplyModeType // enum lora_apply_mode_t lora_apply_mode;
	OffloadParamsToCPU          uint8             // bool offload_params_to_cpu;
	EnableMMAP                  uint8             // bool enable_mmap;
	KeepClipOnCPU               uint8             // bool keep_clip_on_cpu;
	KeepControlNetOnCPU         uint8             // bool keep_control_net_on_cpu;
	KeepVAEOnCPU                uint8             // bool keep_vae_on_cpu;
	FlashAttn                   uint8             // bool flash_attn;
	DiffusionFlashAttn          uint8             // bool diffusion_flash_attn;
	TAEPreviewOnly              uint8             // bool tae_preview_only;
	DiffusionConvDirect         uint8             // bool diffusion_conv_direct;
	VAEConvDirect               uint8             // bool vae_conv_direct;
	CircularX                   uint8             // bool circular_x;
	CircularY                   uint8             // bool circular_y;
	ForceSDXLVAEConvScale       uint8             // bool force_sdxl_vae_conv_scale;
	ChromaUseDITMask            uint8             // bool chroma_use_dit_mask;
	ChromaUseT5Mask             uint8             // bool chroma_use_t5_mask;
	ChromaT5MaskPad             int32             // int chroma_t5_mask_pad;
	QwenImageZeroCond           uint8             // bool qwen_image_zero_cond_t;
}

func (ctx *contextParams) toGo() *ContextParams {
	var _embedding *Embedding
	if ctx.EmbeddingCount != 0 {
		_embedding = ctx.Embeddings.toGo()
	}
	return &ContextParams{
		ModelPath:                   charToString(ctx.ModelPath),
		ClipLPath:                   charToString(ctx.ClipLPath),
		ClipGPath:                   charToString(ctx.ClipGPath),
		ClipVisionPath:              charToString(ctx.ClipVisionPath),
		T5XXLPath:                   charToString(ctx.T5XXLPath),
		LLMPath:                     charToString(ctx.LLMPath),
		LLMVisionPath:               charToString(ctx.LLMVisionPath),
		DiffusionModelPath:          charToString(ctx.DiffusionModelPath),
		HighNoiseDiffusionModelPath: charToString(ctx.HighNoiseDiffusionModelPath),
		VAEPath:                     charToString(ctx.VAEPath),
		TAESDPath:                   charToString(ctx.TAESDPath),
		ControlNetPath:              charToString(ctx.ControlNetPath),
		Embeddings:                  _embedding,
		EmbeddingCount:              ctx.EmbeddingCount,
		PhotoMakerPath:              charToString(ctx.PhotoMakerPath),
		TensorTypeRules:             charToString(ctx.TensorTypeRules),
		VAEDecodeOnly:               byteToBool(ctx.VAEDecodeOnly),
		FreeParamsImmediately:       byteToBool(ctx.FreeParamsImmediately),
		NThreads:                    ctx.NThreads,
		WType:                       ctx.WType,
		RNG:                         ctx.RNG,
		SamplerRNG:                  ctx.SamplerRNG,
		Prediction:                  ctx.Prediction,
		LoraApplyMode:               ctx.LoraApplyMode,
		OffloadParamsToCPU:          byteToBool(ctx.OffloadParamsToCPU),
		EnableMMAP:                  byteToBool(ctx.EnableMMAP),
		KeepClipOnCPU:               byteToBool(ctx.KeepClipOnCPU),
		KeepControlNetOnCPU:         byteToBool(ctx.KeepControlNetOnCPU),
		KeepVAEOnCPU:                byteToBool(ctx.KeepVAEOnCPU),
		FlashAttn:                   byteToBool(ctx.FlashAttn),
		DiffusionFlashAttn:          byteToBool(ctx.DiffusionFlashAttn),
		TAEPreviewOnly:              byteToBool(ctx.TAEPreviewOnly),
		DiffusionConvDirect:         byteToBool(ctx.DiffusionConvDirect),
		VAEConvDirect:               byteToBool(ctx.VAEConvDirect),
		CircularX:                   byteToBool(ctx.CircularX),
		CircularY:                   byteToBool(ctx.CircularY),
		ForceSDXLVAEConvScale:       byteToBool(ctx.ForceSDXLVAEConvScale),
		ChromaUseDITMask:            byteToBool(ctx.ChromaUseDITMask),
		ChromaUseT5Mask:             byteToBool(ctx.ChromaUseT5Mask),
		ChromaT5MaskPad:             ctx.ChromaT5MaskPad,
		QwenImageZeroCond:           byteToBool(ctx.QwenImageZeroCond),
	}
}

type ContextParams struct {
	ModelPath                   string
	ClipLPath                   string
	ClipGPath                   string
	ClipVisionPath              string
	T5XXLPath                   string
	LLMPath                     string
	LLMVisionPath               string
	DiffusionModelPath          string
	HighNoiseDiffusionModelPath string
	VAEPath                     string
	TAESDPath                   string
	ControlNetPath              string
	Embeddings                  *Embedding
	EmbeddingCount              uint32
	PhotoMakerPath              string
	TensorTypeRules             string
	VAEDecodeOnly               bool
	FreeParamsImmediately       bool
	NThreads                    int32
	WType                       SDType
	RNG                         RNGType
	SamplerRNG                  RNGType
	Prediction                  PredictionType
	LoraApplyMode               LoraApplyModeType
	OffloadParamsToCPU          bool
	EnableMMAP                  bool
	KeepClipOnCPU               bool
	KeepControlNetOnCPU         bool
	KeepVAEOnCPU                bool
	FlashAttn                   bool
	DiffusionFlashAttn          bool
	TAEPreviewOnly              bool
	DiffusionConvDirect         bool
	VAEConvDirect               bool
	CircularX                   bool
	CircularY                   bool
	ForceSDXLVAEConvScale       bool
	ChromaUseDITMask            bool
	ChromaUseT5Mask             bool
	ChromaT5MaskPad             int32
	QwenImageZeroCond           bool
}

func (ctx *ContextParams) toC() *contextParams {
	var _embedding *embedding
	if ctx.EmbeddingCount != 0 {
		_embedding = ctx.Embeddings.toC()
	}
	return &contextParams{
		ModelPath:                   stringToChar(ctx.ModelPath),
		ClipLPath:                   stringToChar(ctx.ClipLPath),
		ClipGPath:                   stringToChar(ctx.ClipGPath),
		ClipVisionPath:              stringToChar(ctx.ClipVisionPath),
		T5XXLPath:                   stringToChar(ctx.T5XXLPath),
		LLMPath:                     stringToChar(ctx.LLMPath),
		LLMVisionPath:               stringToChar(ctx.LLMVisionPath),
		DiffusionModelPath:          stringToChar(ctx.DiffusionModelPath),
		HighNoiseDiffusionModelPath: stringToChar(ctx.HighNoiseDiffusionModelPath),
		VAEPath:                     stringToChar(ctx.VAEPath),
		TAESDPath:                   stringToChar(ctx.TAESDPath),
		ControlNetPath:              stringToChar(ctx.ControlNetPath),
		Embeddings:                  _embedding,
		EmbeddingCount:              ctx.EmbeddingCount,
		PhotoMakerPath:              stringToChar(ctx.PhotoMakerPath),
		TensorTypeRules:             stringToChar(ctx.TensorTypeRules),
		VAEDecodeOnly:               boolToByte(ctx.VAEDecodeOnly),
		FreeParamsImmediately:       boolToByte(ctx.FreeParamsImmediately),
		NThreads:                    ctx.NThreads,
		WType:                       ctx.WType,
		RNG:                         ctx.RNG,
		SamplerRNG:                  ctx.SamplerRNG,
		Prediction:                  ctx.Prediction,
		LoraApplyMode:               ctx.LoraApplyMode,
		OffloadParamsToCPU:          boolToByte(ctx.OffloadParamsToCPU),
		EnableMMAP:                  boolToByte(ctx.EnableMMAP),
		KeepClipOnCPU:               boolToByte(ctx.KeepClipOnCPU),
		KeepControlNetOnCPU:         boolToByte(ctx.KeepControlNetOnCPU),
		KeepVAEOnCPU:                boolToByte(ctx.KeepVAEOnCPU),
		FlashAttn:                   boolToByte(ctx.FlashAttn),
		DiffusionFlashAttn:          boolToByte(ctx.DiffusionFlashAttn),
		TAEPreviewOnly:              boolToByte(ctx.TAEPreviewOnly),
		DiffusionConvDirect:         boolToByte(ctx.DiffusionConvDirect),
		VAEConvDirect:               boolToByte(ctx.VAEConvDirect),
		CircularX:                   boolToByte(ctx.CircularX),
		CircularY:                   boolToByte(ctx.CircularY),
		ForceSDXLVAEConvScale:       boolToByte(ctx.ForceSDXLVAEConvScale),
		ChromaUseDITMask:            boolToByte(ctx.ChromaUseDITMask),
		ChromaUseT5Mask:             boolToByte(ctx.ChromaUseT5Mask),
		ChromaT5MaskPad:             ctx.ChromaT5MaskPad,
		QwenImageZeroCond:           boolToByte(ctx.QwenImageZeroCond),
	}
}

type LoraApplyModeType int32

const (
	LoraApplyAuto LoraApplyModeType = iota
	LoraApplyImmediately
	LoraApplyAtRuntime
	LoraApplyModeCount
)

type PredictionType int32

const (
	EPSPred PredictionType = iota
	VPred
	EDMVPred
	FLOWPred
	FLUXFLOWPred
	FLUX2FLOWPred
	PredictionCount
)

type RNGType int32

const (
	STDDefaultRNG RNGType = iota
	CUDARNG
	CPURNG
	RNGTypeCount
)

// enum ggml_type
type SDType int32

//nolint:staticcheck // SA9004: matches C++ enum definition exactly
const (
	TypeF32  SDType = 0
	TypeF16         = 1
	TypeQ4_0        = 2
	TypeQ4_1        = 3
	// SD_TYPE_Q4_2 = 4, support has been removed
	// SD_TYPE_Q4_3 = 5, support has been removed
	TypeQ5_0    = 6
	TypeQ5_1    = 7
	TypeQ8_0    = 8
	TypeQ8_1    = 9
	TypeQ2_K    = 10
	TypeQ3_K    = 11
	TypeQ4_K    = 12
	TypeQ5_K    = 13
	TypeQ6_K    = 14
	TypeQ8_K    = 15
	TypeIQ2_XXS = 16
	TypeIQ2_XS  = 17
	TypeIQ3_XXS = 18
	TypeIQ1_S   = 19
	TypeIQ4_NL  = 20
	TypeIQ3_S   = 21
	TypeIQ2_S   = 22
	TypeIQ4_XS  = 23
	TypeI8      = 24
	TypeI16     = 25
	TypeI32     = 26
	TypeI64     = 27
	TypeF64     = 28
	TypeIQ1_M   = 29
	TypeBF16    = 30
	// SD_TYPE_Q4_0_4_4 = 31, support has been removed from gguf files
	// SD_TYPE_Q4_0_4_8 = 32,
	// SD_TYPE_Q4_0_8_8 = 33,
	TypeTQ1_0 = 34
	TypeTQ2_0 = 35
	// SD_TYPE_IQ4_NL_4_4 = 36,
	// SD_TYPE_IQ4_NL_4_8 = 37,
	// SD_TYPE_IQ4_NL_8_8 = 38,
	TypeMXFP4 = 39 // MXFP4 (1 block)
	TypeCOUNT = 40
)

type embedding struct {
	Name *byte
	Path *byte
}

func (e *embedding) toGo() *Embedding {
	return &Embedding{
		Name: charToString(e.Name),
		Path: charToString(e.Path),
	}
}

type Embedding struct {
	Name string
	Path string
}

func (e *Embedding) toC() *embedding {
	return &embedding{
		Name: stringToChar(e.Name),
		Path: stringToChar(e.Path),
	}
}

// Creates default context params
func ContextParamsInit() ContextParams {
	cp := newContextParams()

	contextParamsInit.Call(nil, unsafe.Pointer(&cp))
	return *cp.toGo()
}

func NewContext(ctxParams ContextParams) Context {
	var context Context

	_ctxParams := ctxParams.toC()
	newContext.Call(unsafe.Pointer(&context), unsafe.Pointer(&_ctxParams))

	return context
}

func FreeCtx(ctx Context) {
	freeCtx.Call(nil, unsafe.Pointer(&ctx))
}

func CtxParamsToStr(ctxParams ContextParams) string {
	str := utilsGetNulString()

	_params := ctxParams.toC()
	ctxParamsToStr.Call(unsafe.Pointer(&str), unsafe.Pointer(&_params))

	return charToString(str)
}

func newContextParams() *contextParams {
	return &contextParams{
		ModelPath:                   utilsGetNulString(),
		ClipLPath:                   utilsGetNulString(),
		ClipGPath:                   utilsGetNulString(),
		ClipVisionPath:              utilsGetNulString(),
		T5XXLPath:                   utilsGetNulString(),
		LLMPath:                     utilsGetNulString(),
		LLMVisionPath:               utilsGetNulString(),
		DiffusionModelPath:          utilsGetNulString(),
		HighNoiseDiffusionModelPath: utilsGetNulString(),
		VAEPath:                     utilsGetNulString(),
		TAESDPath:                   utilsGetNulString(),
		ControlNetPath:              utilsGetNulString(),
		PhotoMakerPath:              utilsGetNulString(),
		TensorTypeRules:             utilsGetNulString(),
	}
}

package gosd

import (
	"unsafe"

	"github.com/jupiterrider/ffi"
)

var (
	// SD_API const char* sd_get_system_info();
	getSystemInfo ffi.Fun

	// SD_API const char* sd_commit(void);
	commit ffi.Fun

	// SD_API const char* sd_version(void);
	version ffi.Fun

	// SD_API int32_t sd_get_num_physical_cores();
	getNumPhysicalCores ffi.Fun

	// SD_API const char* sd_type_name(enum sd_type_t type);
	typeName ffi.Fun

	// SD_API enum sd_type_t str_to_sd_type(const char* str);
	strToSDType ffi.Fun

	// SD_API const char* sd_rng_type_name(enum rng_type_t rng_type);
	rngTypeName ffi.Fun

	// SD_API enum rng_type_t str_to_rng_type(const char* str);
	strToRNGType ffi.Fun

	// SD_API const char* sd_sample_method_name(enum sample_method_t sample_method);
	sampleMethodName ffi.Fun

	// SD_API enum sample_method_t str_to_sample_method(const char* str);
	strToSampleMethod ffi.Fun

	// SD_API const char* sd_scheduler_name(enum scheduler_t scheduler);
	schedulerName ffi.Fun

	// SD_API enum scheduler_t str_to_scheduler(const char* str);
	strToScheduler ffi.Fun

	// SD_API const char* sd_prediction_name(enum prediction_t prediction);
	predictionName ffi.Fun

	// SD_API enum prediction_t str_to_prediction(const char* str);
	strToPrediction ffi.Fun

	// SD_API const char* sd_preview_name(enum preview_t preview);
	previewName ffi.Fun

	// SD_API enum preview_t str_to_preview(const char* str);
	strToPreview ffi.Fun

	// SD_API const char* sd_lora_apply_mode_name(enum lora_apply_mode_t mode);
	loraApplyModeName ffi.Fun

	// SD_API enum lora_apply_mode_t str_to_lora_apply_mode(const char* str);
	strToLoraApplyMode ffi.Fun

	// SD_API const char* sd_hires_upscaler_name(enum sd_hires_upscaler_t upscaler);
	hiresUpscalerName ffi.Fun

	// SD_API enum sd_hires_upscaler_t str_to_sd_hires_upscaler(const char* str);
	strToHiresUpscaler ffi.Fun

	// SD_API bool convert(const char* input_path, const char* vae_path, const char* output_path, enum sd_type_t output_type, const char* tensor_type_rules, bool convert_name);
	convert ffi.Fun

	// SD_API bool preprocess_canny(sd_image_t image, float high_threshold, float low_threshold, float weak, float strong, bool inverse);
	preprocessCanny ffi.Fun
)

func loadSystemRoutines(lib ffi.Lib) error {
	var err error
	if getSystemInfo, err = lib.Prep("sd_get_system_info", &ffi.TypePointer); err != nil {
		return loadError("sd_get_system_info", err)
	}

	if commit, err = lib.Prep("sd_commit", &ffi.TypePointer); err != nil {
		return loadError("sd_commit", err)
	}

	if version, err = lib.Prep("sd_version", &ffi.TypePointer); err != nil {
		return loadError("sd_version", err)
	}

	if getNumPhysicalCores, err = lib.Prep("sd_get_num_physical_cores", &ffi.TypeSint32); err != nil {
		return loadError("sd_get_num_physical_cores", err)
	}

	if typeName, err = lib.Prep("sd_type_name", &ffi.TypePointer, &ffi.TypeSint32); err != nil {
		return loadError("sd_type_name", err)
	}

	if strToSDType, err = lib.Prep("str_to_sd_type", &ffi.TypeSint32, &ffi.TypePointer); err != nil {
		return loadError("str_to_sd_type", err)
	}

	if rngTypeName, err = lib.Prep("sd_rng_type_name", &ffi.TypePointer, &ffi.TypeSint32); err != nil {
		return loadError("sd_rng_type_name", err)
	}

	if strToRNGType, err = lib.Prep("str_to_rng_type", &ffi.TypeSint32, &ffi.TypePointer); err != nil {
		return loadError("str_to_rng_type", err)
	}

	if sampleMethodName, err = lib.Prep("sd_sample_method_name", &ffi.TypePointer, &ffi.TypeSint32); err != nil {
		return loadError("sd_sample_method_name", err)
	}

	if strToSampleMethod, err = lib.Prep("str_to_sample_method", &ffi.TypeSint32, &ffi.TypePointer); err != nil {
		return loadError("str_to_sample_method", err)
	}

	if schedulerName, err = lib.Prep("sd_scheduler_name", &ffi.TypePointer, &ffi.TypeSint32); err != nil {
		return loadError("sd_scheduler_name", err)
	}

	if strToScheduler, err = lib.Prep("str_to_scheduler", &ffi.TypeSint32, &ffi.TypePointer); err != nil {
		return loadError("str_to_scheduler", err)
	}

	if predictionName, err = lib.Prep("sd_prediction_name", &ffi.TypePointer, &ffi.TypeSint32); err != nil {
		return loadError("sd_prediction_name", err)
	}

	if strToPrediction, err = lib.Prep("str_to_prediction", &ffi.TypeSint32, &ffi.TypePointer); err != nil {
		return loadError("str_to_prediction", err)
	}

	if previewName, err = lib.Prep("sd_preview_name", &ffi.TypePointer, &ffi.TypeSint32); err != nil {
		return loadError("sd_preview_name", err)
	}

	if strToPreview, err = lib.Prep("str_to_preview", &ffi.TypeSint32, &ffi.TypePointer); err != nil {
		return loadError("str_to_preview", err)
	}

	if loraApplyModeName, err = lib.Prep("sd_lora_apply_mode_name", &ffi.TypePointer, &ffi.TypeSint32); err != nil {
		return loadError("sd_lora_apply_mode_name", err)
	}

	if strToLoraApplyMode, err = lib.Prep("str_to_lora_apply_mode", &ffi.TypeSint32, &ffi.TypePointer); err != nil {
		return loadError("str_to_lora_apply_mode", err)
	}

	if hiresUpscalerName, err = lib.Prep("sd_hires_upscaler_name", &ffi.TypePointer, &ffi.TypeSint32); err != nil {
		return loadError("sd_hires_upscaler_name", err)
	}

	if strToHiresUpscaler, err = lib.Prep("str_to_sd_hires_upscaler", &ffi.TypeSint32, &ffi.TypePointer); err != nil {
		return loadError("str_to_sd_hires_upscaler", err)
	}

	if convert, err = lib.Prep("convert", &ffi.TypeUint8, &ffi.TypePointer, &ffi.TypePointer, &ffi.TypePointer, &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypeUint8); err != nil {
		return loadError("convert", err)
	}

	if preprocessCanny, err = lib.Prep("preprocess_canny", &ffi.TypeUint8, &ffiTypeImage, &ffi.TypeFloat, &ffi.TypeFloat, &ffi.TypeFloat, &ffi.TypeFloat, &ffi.TypeUint8); err != nil {
		return loadError("preprocess_canny", err)
	}

	return nil
}

// GetSystemInfo returns a formatted string containing the CPU instruction sets
// supported by the current hardware (e.g., AVX, AVX2, FMA). This is used
// to verify which hardware acceleration features are active for GGML operations.
func GetSystemInfo() string {
	var systemInfo *byte

	getSystemInfo.Call(unsafe.Pointer(&systemInfo))
	if systemInfo == nil {
		return ""
	}

	return charToString(systemInfo)
}

// Commit returns stable-diffusion.cpp commit hash.
func Commit() string {
	var commitInfo *byte

	commit.Call(unsafe.Pointer(&commitInfo))
	if commitInfo == nil {
		return ""
	}

	return charToString(commitInfo)
}

// Version returns stable-diffusion.cpp release version.
func Version() string {
	var versionInfo *byte

	version.Call(unsafe.Pointer(&versionInfo))
	if versionInfo == nil {
		return ""
	}

	return charToString(versionInfo)
}

// GetNumPhysicalCores returns number of physical cores in the system.
func GetNumPhysicalCores() int {
	var count int

	getNumPhysicalCores.Call(unsafe.Pointer(&count))
	return count
}

// TypeName stringifies SDType
func TypeName(sdType SDType) string {
	res := utilsGetNulString()

	typeName.Call(unsafe.Pointer(&res), unsafe.Pointer(&sdType))
	return charToString(res)
}

// StrToSDType converts SDType name to its enumeration.
func StrToSDType(typeName string) SDType {
	var sdType SDType
	name := utilsStrToNulString(typeName)

	strToSDType.Call(unsafe.Pointer(&sdType), unsafe.Pointer(&name))
	return sdType
}

// RNGTypeName stringifies RNGType.
func RNGTypeName(rngType RNGType) string {
	res := utilsGetNulString()

	rngTypeName.Call(unsafe.Pointer(&res), unsafe.Pointer(&rngType))
	return charToString(res)
}

// StrToRNGType converts RNGType name to its enumeration.
func StrToRNGType(typeName string) RNGType {
	var rngType RNGType
	name := utilsStrToNulString(typeName)

	strToRNGType.Call(unsafe.Pointer(&rngType), unsafe.Pointer(&name))
	return rngType
}

// SampleMethodName stringifies SampleMethodType.
func SampleMethodName(sampleMethod SampleMethodType) string {
	res := utilsGetNulString()

	sampleMethodName.Call(unsafe.Pointer(&res), unsafe.Pointer(&sampleMethod))
	return charToString(res)
}

// StrToSampleMethod converts SampleMethodType name to its enumeration.
func StrToSampleMethod(typeName string) SampleMethodType {
	var sampleMethodType SampleMethodType
	name := utilsStrToNulString(typeName)

	strToSampleMethod.Call(unsafe.Pointer(&sampleMethodType), unsafe.Pointer(&name))
	return sampleMethodType
}

// SchedulerName stringifies SchedulerType.
func SchedulerName(schedulerType SchedulerType) string {
	res := utilsGetNulString()

	schedulerName.Call(unsafe.Pointer(&res), unsafe.Pointer(&schedulerType))
	return charToString(res)
}

// StrToScheduler converts SchedulerType name to its enumeration.
func StrToScheduler(typeName string) SchedulerType {
	var schedulerType SchedulerType
	name := utilsStrToNulString(typeName)

	strToScheduler.Call(unsafe.Pointer(&schedulerType), unsafe.Pointer(&name))
	return schedulerType
}

// PredictionName stringifies PredictionType.
func PredictionName(predictionType PredictionType) string {
	res := utilsGetNulString()

	predictionName.Call(unsafe.Pointer(&res), unsafe.Pointer(&predictionType))
	return charToString(res)
}

// StrToPrediction converts PredictionType name to its enumeration.
func StrToPrediction(typeName string) PredictionType {
	var predictionType PredictionType
	name := utilsStrToNulString(typeName)

	strToPrediction.Call(unsafe.Pointer(&predictionType), unsafe.Pointer(&name))
	return predictionType
}

// PreviewName stringifies PreviewMode
func PreviewName(previewType PreviewMode) string {
	res := utilsGetNulString()

	previewName.Call(unsafe.Pointer(&res), unsafe.Pointer(&previewType))
	return charToString(res)
}

// StrToPreview converts PreviewMode name to its enumeration.
func StrToPreview(typeName string) PreviewMode {
	var previewMode PreviewMode
	name := utilsStrToNulString(typeName)

	strToPreview.Call(unsafe.Pointer(&previewMode), unsafe.Pointer(&name))
	return previewMode
}

// LoraApplyModeName stringifies LoraApplyModeType.
func LoraApplyModeName(loraMode LoraApplyModeType) string {
	res := utilsGetNulString()

	loraApplyModeName.Call(unsafe.Pointer(&res), unsafe.Pointer(&loraMode))
	return charToString(res)
}

// StrToLoraApplyMode converts LoraApplyModeType to its enumeration.
func StrToLoraApplyMode(typeName string) LoraApplyModeType {
	var loraMode LoraApplyModeType
	name := utilsStrToNulString(typeName)

	strToLoraApplyMode.Call(unsafe.Pointer(&loraMode), unsafe.Pointer(&name))
	return loraMode
}

// HiresUpscalerName stringifies HiresUpscalerType.
func HiresUpscalerName(hiresMode HiresUpscalerType) string {
	res := utilsGetNulString()

	hiresUpscalerName.Call(unsafe.Pointer(&res), unsafe.Pointer(&hiresMode))
	return charToString(res)
}

// StrToHiresUpscaler converts HiresUpscalerType to its enumeration.
func StrToHiresUpscaler(typeName string) HiresUpscalerType {
	var hiresMode HiresUpscalerType
	name := utilsStrToNulString(typeName)

	strToHiresUpscaler.Call(unsafe.Pointer(&hiresMode), unsafe.Pointer(&name))
	return hiresMode
}

// Convert converts model to safetensor/gguf format.
// If VAE model is provided it will be merged with the diffusion model.
// CPU-bound API.
func Convert(modelPath string, vaePath string, outputPath string, outputType SDType, tensorTypeRules string, convertName bool) bool {
	mp := stringToChar(modelPath)
	vp := stringToChar(vaePath)
	op := stringToChar(outputPath)
	ttr := stringToChar(tensorTypeRules)
	cn := boolToByte(convertName)

	res := uint8(0)

	convert.Call(
		unsafe.Pointer(&res),
		unsafe.Pointer(&mp),
		unsafe.Pointer(&vp),
		unsafe.Pointer(&op),
		unsafe.Pointer(&outputType),
		unsafe.Pointer(&ttr),
		unsafe.Pointer(&cn),
	)
	return byteToBool(res)
}

// PreprocessCanny applies Canny algorithm for edge detection in an image.
// CPU-bound API.
func PreprocessCanny(image Image, highThreshold float32, lowThreshold float32, weak float32, strong float32, inverse bool) bool {
	img := *image.toC()
	inv := boolToByte(inverse)

	res := uint8(0)

	preprocessCanny.Call(
		unsafe.Pointer(&res),
		unsafe.Pointer(&img),
		unsafe.Pointer(&highThreshold),
		unsafe.Pointer(&lowThreshold),
		unsafe.Pointer(&weak),
		unsafe.Pointer(&strong),
		unsafe.Pointer(&inv),
	)
	return byteToBool(res)
}

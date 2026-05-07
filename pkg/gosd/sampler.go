package gosd

import (
	"unsafe"

	"github.com/jupiterrider/ffi"
)

var (
	// SD_API void sd_sample_params_init(sd_sample_params_t* sample_params);
	sampleParamsInit ffi.Fun

	// SD_API char* sd_sample_params_to_str(const sd_sample_params_t* sample_params);
	sampleParamsToStr ffi.Fun

	// SD_API enum sample_method_t sd_get_default_sample_method(const sd_ctx_t* sd_ctx);
	getDefaultSampleMethod ffi.Fun

	// SD_API enum scheduler_t sd_get_default_scheduler(const sd_ctx_t* sd_ctx, enum sample_method_t sample_method);
	getDefaultScheduler ffi.Fun
)

func loadSamplerRoutines(lib ffi.Lib) error {
	var err error

	if sampleParamsInit, err = lib.Prep(
		"sd_sample_params_init",
		&ffi.TypeVoid,
		&ffi.TypePointer,
	); err != nil {
		return loadError("sd_sample_params_init", err)
	}

	if sampleParamsToStr, err = lib.Prep(
		"sd_sample_params_to_str",
		&ffi.TypePointer,
		&ffi.TypePointer,
	); err != nil {
		return loadError("sd_sample_params_to_str", err)
	}

	if getDefaultSampleMethod, err = lib.Prep(
		"sd_get_default_sample_method",
		&ffi.TypeSint32,
		&ffi.TypePointer,
	); err != nil {
		return loadError("sd_get_default_sample_method", err)
	}

	if getDefaultScheduler, err = lib.Prep(
		"sd_get_default_scheduler",
		&ffi.TypeSint32,
		&ffi.TypePointer,
		&ffi.TypeSint32,
	); err != nil {
		return loadError("sd_get_default_scheduler", err)
	}

	return nil
}

type sampleParamsType struct {
	Guidance          guidanceParams   // sd_guidance_params_t guidance;
	Scheduler         SchedulerType    // enum scheduler_t scheduler;
	SampleMethod      SampleMethodType // enum sample_method_t sample_method;
	SampleSteps       int32            // int sample_steps;
	ETA               float32          // float eta;
	ShiftedTimestamp  int32            // int shifted_timestep;
	CustomSigmas      *float32         // float* custom_sigmas;
	CustomSigmasCount int32            // int custom_sigmas_count;
	FlowShift         float32          // float flow_shift;
}

func (slg *sampleParamsType) toGo() *SampleParamsType {
	size := int(slg.CustomSigmasCount)
	newSigma := make([]float32, size)

	srcSigma := unsafe.Slice(slg.CustomSigmas, size)
	copy(newSigma, srcSigma)

	return &SampleParamsType{
		Guidance:          *slg.Guidance.toGo(),
		Scheduler:         slg.Scheduler,
		SampleMethod:      slg.SampleMethod,
		SampleSteps:       slg.SampleSteps,
		ETA:               slg.ETA,
		ShiftedTimestamp:  slg.ShiftedTimestamp,
		CustomSigmas:      newSigma,
		CustomSigmasCount: slg.CustomSigmasCount,
		FlowShift:         slg.FlowShift,
	}
}

type SampleParamsType struct {
	Guidance          GuidanceParams
	Scheduler         SchedulerType
	SampleMethod      SampleMethodType
	SampleSteps       int32
	ETA               float32
	ShiftedTimestamp  int32
	CustomSigmas      []float32
	CustomSigmasCount int32
	FlowShift         float32
}

func (slg *SampleParamsType) toC() *sampleParamsType {
	size := int(slg.CustomSigmasCount)
	var _data *float32

	if size != 0 {
		_data = &slg.CustomSigmas[0]
	}

	return &sampleParamsType{
		Guidance:          *slg.Guidance.toC(),
		Scheduler:         slg.Scheduler,
		SampleMethod:      slg.SampleMethod,
		SampleSteps:       slg.SampleSteps,
		ETA:               slg.ETA,
		ShiftedTimestamp:  slg.ShiftedTimestamp,
		CustomSigmas:      _data,
		CustomSigmasCount: slg.CustomSigmasCount,
		FlowShift:         slg.FlowShift,
	}
}

func newSampleParams() *sampleParamsType {
	return &sampleParamsType{}
}

// SampleParamsInit initializes default values for the inference sampler.
func SampleParamsInit() SampleParamsType {
	sp := newSampleParams()

	sampleParamsInit.Call(nil, unsafe.Pointer(&sp))

	return *sp.toGo()
}

// SampleParamsToStr stringifies structure encapsulating sampler parameters.
func SampleParamsToStr(params SampleParamsType) string {
	sp := params.toC()
	str := utilsGetNulString()

	sampleParamsToStr.Call(unsafe.Pointer(&str), unsafe.Pointer(&sp))

	return charToString(str)
}

// GetDefaultSampleMethod returns default sampler method from a context.
func GetDefaultSampleMethod(ctx Context) SampleMethodType {
	var sampleType SampleMethodType

	getDefaultSampleMethod.Call(unsafe.Pointer(&sampleType), unsafe.Pointer(&ctx))

	return sampleType
}

// GetDefaultScheduler returns default scheduler type from a context.
func GetDefaultScheduler(ctx Context, sampler SampleMethodType) SchedulerType {
	var schedulerType SchedulerType

	getDefaultScheduler.Call(unsafe.Pointer(&schedulerType), unsafe.Pointer(&ctx), unsafe.Pointer(&sampler))

	return schedulerType
}

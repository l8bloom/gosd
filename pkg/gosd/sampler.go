package gosd

import "unsafe"

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

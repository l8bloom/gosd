package gosd

import (
	"unsafe"

	"github.com/jupiterrider/ffi"
)

var (
	// SD_API void sd_cache_params_init(sd_cache_params_t* cache_params);  // +
	cacheParamsInit ffi.Fun
)

func loadCacheRoutines(lib ffi.Lib) error {
	var err error

	if cacheParamsInit, err = lib.Prep(
		"sd_cache_params_init",
		&ffi.TypeVoid,
		&ffi.TypePointer,
	); err != nil {
		return loadError("sd_cache_params_init", err)
	}

	return nil
}

type CacheModeType int32

const (
	CacheDisabled CacheModeType = iota
	CacheEasyCache
	CacheUcache
	CacheDBcache
	CacheTaylorseer
	CacheCacheDit
	CacheSpectrum
)

type cacheParams struct {
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
	SpectrumW                float32       // float spectrum_w;
	SpectrumM                int32         // int spectrum_m;
	SpectrumLAM              float32       // float spectrum_lam;
	SpectrumWindowSize       int32         // int spectrum_window_size;
	SpectrumFlexWindow       float32       // float spectrum_flex_window;
	SpectrumWarmupSteps      int32         // int spectrum_warmup_steps;
	SpectrumStopPercent      float32       // float spectrum_stop_percent;
}

func (c *cacheParams) toGo() *CacheParams {
	return &CacheParams{
		Mode:                     c.Mode,
		ReuseThreshold:           c.ReuseThreshold,
		StartPercent:             c.StartPercent,
		EndPercent:               c.EndPercent,
		ErrorDecayRate:           c.ErrorDecayRate,
		UseRelativeThreshold:     byteToBool(c.UseRelativeThreshold),
		ResetErrorOnCompute:      byteToBool(c.ResetErrorOnCompute),
		FNComputeBlocks:          c.FNComputeBlocks,
		BNComputeBlocks:          c.BNComputeBlocks,
		ResidualDiffThreshold:    c.ResidualDiffThreshold,
		MaxWarmupSteps:           c.MaxWarmupSteps,
		MaxCachedSteps:           c.MaxCachedSteps,
		MaxContinuousCachedSteps: c.MaxContinuousCachedSteps,
		TaylorSeerNDerivatives:   c.TaylorSeerNDerivatives,
		TaylorSeerSkipInterval:   c.TaylorSeerSkipInterval,
		SCMMask:                  charToString(c.SCMMask),
		SCMPolicyDynamic:         byteToBool(c.SCMPolicyDynamic),
		SpectrumW:                c.SpectrumW,
		SpectrumM:                c.SpectrumM,
		SpectrumLAM:              c.SpectrumLAM,
		SpectrumWindowSize:       c.SpectrumWindowSize,
		SpectrumFlexWindow:       c.SpectrumFlexWindow,
		SpectrumWarmupSteps:      c.SpectrumWarmupSteps,
		SpectrumStopPercent:      c.SpectrumStopPercent,
	}
}

type CacheParams struct {
	Mode                     CacheModeType
	ReuseThreshold           float32
	StartPercent             float32
	EndPercent               float32
	ErrorDecayRate           float32
	UseRelativeThreshold     bool
	ResetErrorOnCompute      bool
	FNComputeBlocks          int32
	BNComputeBlocks          int32
	ResidualDiffThreshold    float32
	MaxWarmupSteps           int32
	MaxCachedSteps           int32
	MaxContinuousCachedSteps int32
	TaylorSeerNDerivatives   int32
	TaylorSeerSkipInterval   int32
	SCMMask                  string
	SCMPolicyDynamic         bool
	SpectrumW                float32
	SpectrumM                int32
	SpectrumLAM              float32
	SpectrumWindowSize       int32
	SpectrumFlexWindow       float32
	SpectrumWarmupSteps      int32
	SpectrumStopPercent      float32
}

func (c *CacheParams) toC() *cacheParams {
	return &cacheParams{
		Mode:                     c.Mode,
		ReuseThreshold:           c.ReuseThreshold,
		StartPercent:             c.StartPercent,
		EndPercent:               c.EndPercent,
		ErrorDecayRate:           c.ErrorDecayRate,
		UseRelativeThreshold:     boolToByte(c.UseRelativeThreshold),
		ResetErrorOnCompute:      boolToByte(c.ResetErrorOnCompute),
		FNComputeBlocks:          c.FNComputeBlocks,
		BNComputeBlocks:          c.BNComputeBlocks,
		ResidualDiffThreshold:    c.ResidualDiffThreshold,
		MaxWarmupSteps:           c.MaxWarmupSteps,
		MaxCachedSteps:           c.MaxCachedSteps,
		MaxContinuousCachedSteps: c.MaxContinuousCachedSteps,
		TaylorSeerNDerivatives:   c.TaylorSeerNDerivatives,
		TaylorSeerSkipInterval:   c.TaylorSeerSkipInterval,
		SCMMask:                  stringToChar(c.SCMMask),
		SCMPolicyDynamic:         boolToByte(c.SCMPolicyDynamic),
		SpectrumW:                c.SpectrumW,
		SpectrumM:                c.SpectrumM,
		SpectrumLAM:              c.SpectrumLAM,
		SpectrumWindowSize:       c.SpectrumWindowSize,
		SpectrumFlexWindow:       c.SpectrumFlexWindow,
		SpectrumWarmupSteps:      c.SpectrumWarmupSteps,
		SpectrumStopPercent:      c.SpectrumStopPercent,
	}
}

func newCacheParams() *cacheParams {
	cp := &cacheParams{
		SCMMask: utilsGetNulString(),
	}

	return cp
}

func CacheParamsInit() CacheParams {
	cp := newCacheParams()

	cacheParamsInit.Call(nil, unsafe.Pointer(&cp))

	return *cp.toGo()
}

package gosd

import (
	"testing"
)

func TestGetSystemInfo(t *testing.T) {
	if err := Load(); err != nil {
		t.Fatal(err.Error())
	}
	systemInfo := GetSystemInfo()

	if len(systemInfo) == 0 {
		t.Errorf("Expected non-empty SystemInfo string, got %s", systemInfo)
	}
}

func TestCommit(t *testing.T) {
	sdCommit := Commit()

	if len(sdCommit) == 0 {
		t.Errorf("Expected non-empty Commit string, got %s", sdCommit)
	}
}

func TestVersion(t *testing.T) {
	version := Version()

	if len(version) == 0 {
		t.Errorf("Expected non-empty Version string, got %s", version)
	}
}

func TestGetNumPhysicalCores(t *testing.T) {
	coresCount := GetNumPhysicalCores()

	if coresCount == 0 {
		t.Errorf("Expected positive number of physical cores, got %d", coresCount)
	}
}

func TestTypeName(t *testing.T) {
	name := TypeName(TypeF32)

	if name != "f32" {
		t.Errorf("expected `f32` for `TypeF32`, got  %q", name)
	}

	name = TypeName(TypeNVFP4)

	if name != "nvfp4" {
		t.Errorf("expected `nvfp4` for `TypeNVFP4`, got  %q", name)
	}
}

func TestStrToSDType(t *testing.T) {
	name := "f32"
	sdType := StrToSDType(name)

	if sdType != 0 {
		t.Errorf("expected `0` for `f32`, got  %d", sdType)
	}

	name = "nvfp4"
	sdType = StrToSDType(name)

	if sdType != 40 {
		t.Errorf("expected `40` for `nvfp4`, got  %d", sdType)
	}
}

func TestRNGTypeName(t *testing.T) {
	name := RNGTypeName(CUDARNG)

	if name != "cuda" {
		t.Errorf("expected `cuda` for `CUDARNG`, got  %q", name)
	}

	name = RNGTypeName(CPURNG)

	if name != "cpu" {
		t.Errorf("expected `cpu` for `CPURNG`, got  %q", name)
	}
}

func TestStrToRNGType(t *testing.T) {
	name := "cuda"
	rngType := StrToRNGType(name)

	if rngType != 1 {
		t.Errorf("expected `1` for `CUDARNG`, got  %d", rngType)
	}

	name = "cpu"
	rngType = StrToRNGType(name)

	if rngType != 2 {
		t.Errorf("expected `2` for `CPURNG`, got  %d", rngType)
	}
}

func TestSampleMethodName(t *testing.T) {
	name := SampleMethodName(EulerSampleMethod)

	if name != "euler" {
		t.Errorf("expected `euler` for `EulerSampleMethod`, got  %q", name)
	}

	name = SampleMethodName(ERSDESampleMethod)

	if name != "er_sde" {
		t.Errorf("expected `er_sde` for `ERSDESampleMethod`, got  %q", name)
	}
}

func TestStrToSampleMethod(t *testing.T) {
	name := "euler"
	sampleType := StrToSampleMethod(name)

	if sampleType != 0 {
		t.Errorf("expected `0` for `EulerSampleMethod`, got  %d", sampleType)
	}

	name = "er_sde"
	sampleType = StrToSampleMethod(name)

	if sampleType != 14 {
		t.Errorf("expected `14` for `ERSDESampleMethod`, got  %d", sampleType)
	}
}

func TestSchedulerName(t *testing.T) {
	name := SchedulerName(DiscreteScheduler)

	if name != "discrete" {
		t.Errorf("expected `discrete` for `DiscreteScheduler`, got  %q", name)
	}

	name = SchedulerName(BongTangentScheduler)

	if name != "bong_tangent" {
		t.Errorf("expected `bong_tangent` for `BongTangentScheduler`, got  %q", name)
	}
}

func TestStrToScheduler(t *testing.T) {
	name := "discrete"
	schedulerType := StrToScheduler(name)

	if schedulerType != 0 {
		t.Errorf("expected `0` for `DiscreteScheduler`, got  %d", schedulerType)
	}

	name = "bong_tangent"
	schedulerType = StrToScheduler(name)

	if schedulerType != 10 {
		t.Errorf("expected `10` for `BongTangentScheduler`, got  %d", schedulerType)
	}
}

func TestPredictionName(t *testing.T) {
	name := PredictionName(EPSPred)

	if name != "eps" {
		t.Errorf("expected `eps` for `EPSPred`, got  %q", name)
	}

	name = PredictionName(FLUX2FLOWPred)

	if name != "flux2_flow" {
		t.Errorf("expected `flux2_flow` for `FLUX2FLOWPred`, got  %q", name)
	}
}

func TestStrToPrediction(t *testing.T) {
	name := "eps"
	predictionType := StrToPrediction(name)

	if predictionType != 0 {
		t.Errorf("expected `0` for `EPSPred`, got  %d", predictionType)
	}

	name = "flux2_flow"
	predictionType = StrToPrediction(name)

	if predictionType != 5 {
		t.Errorf("expected `5` for `FLUX2FLOWPred`, got  %d", predictionType)
	}
}

func TestPreviewName(t *testing.T) {
	name := PreviewName(PreviewPROJ)

	if name != "proj" {
		t.Errorf("expected `proj` for `PreviewPROJ`, got  %q", name)
	}

	name = PreviewName(PreviewVAE)

	if name != "vae" {
		t.Errorf("expected `vae` for `PreviewVAE`, got  %q", name)
	}
}

func TestStrToPreview(t *testing.T) {
	name := "proj"
	previewMode := StrToPreview(name)

	if previewMode != 1 {
		t.Errorf("expected `1` for `PreviewPROJ`, got  %d", previewMode)
	}

	name = "vae"
	previewMode = StrToPreview(name)

	if previewMode != 3 {
		t.Errorf("expected `3` for `PreviewVAE`, got  %d", previewMode)
	}
}

func TestLoraApplyModeName(t *testing.T) {
	name := LoraApplyModeName(LoraApplyImmediately)

	if name != "immediately" {
		t.Errorf("expected `immediately` for `LoraApplyImmediately`, got  %q", name)
	}

	name = LoraApplyModeName(LoraApplyAtRuntime)

	if name != "at_runtime" {
		t.Errorf("expected `at_runtime` for `LoraApplyAtRuntime`, got  %q", name)
	}
}

func TestStrToLoraApplyMode(t *testing.T) {
	name := "immediately"
	loraMode := StrToLoraApplyMode(name)

	if loraMode != 1 {
		t.Errorf("expected `1` for `LoraApplyImmediately`, got  %d", loraMode)
	}

	name = "at_runtime"
	loraMode = StrToLoraApplyMode(name)

	if loraMode != 2 {
		t.Errorf("expected `2` for `LoraApplyAtRuntime`, got  %d", loraMode)
	}
}

func TestHiresUpscalerName(t *testing.T) {
	name := HiresUpscalerName(HiresUpscalerNone)

	if name != "None" {
		t.Errorf("expected `None` for `HiresUpscalerNone`, got  %q", name)
	}

	name = HiresUpscalerName(HiresUpscalerLatent)

	if name != "Latent" {
		t.Errorf("expected `Latent` for `HiresUpscalerLatent`, got  %q", name)
	}

	name = HiresUpscalerName(HiresUpscalerModel)

	if name != "Model" {
		t.Errorf("expected `Model` for `HiresUpscalerModel`, got  %q", name)
	}

}

func TestStrToHiresUpscaler(t *testing.T) {
	name := "Latent (nearest-exact)"
	hiresMode := StrToHiresUpscaler(name)

	if hiresMode != HiresUpscalerLatentNearestExact {
		t.Errorf("expected `%d` for `nearest-exact`, got  %d", HiresUpscalerLatentNearestExact, hiresMode)
	}

	name = "Latent (bicubic antialiased)"
	hiresMode = StrToHiresUpscaler(name)

	if hiresMode != HiresUpscalerLatentBicubicAntialiased {
		t.Errorf("expected `%d` for `bicubic antialiased`, got  %d", HiresUpscalerLatentBicubic, hiresMode)
	}
}

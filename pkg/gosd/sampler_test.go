package gosd

import (
	"testing"
)

func TestSampleParamsInit(t *testing.T) {
	Load()
	sp := SampleParamsInit()

	if sp.SampleSteps != 20 {
		t.Errorf("expected SampleSteps=20, got %d", sp.SampleSteps)
	}

	if sp.Guidance.TextCfg != 7 {
		t.Errorf("expected Guidance.TextCfg=7, got %g", sp.Guidance.TextCfg)
	}

	if sp.Guidance.DistilledGuidance != 3.5 {
		t.Errorf("expected Guidance.DistilledGuidance=3.5, got %g", sp.Guidance.DistilledGuidance)
	}
}

func TestSampleParamsToStr(t *testing.T) {
	sp := SampleParamsInit()
	spStr := SampleParamsToStr(sp)

	if len(spStr) == 0 {
		t.Errorf("expected non-empty sampler params string, got %s", spStr)
	}
}

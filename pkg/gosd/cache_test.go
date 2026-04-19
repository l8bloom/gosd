package gosd

import (
	"testing"
)

// test only some sensible default values
func TestCacheParamsInit(t *testing.T) {
	Load()
	cp := CacheParamsInit()

	if cp.Mode != CacheDisabled {
		t.Errorf("expected Mode=CacheDisabled, got %d", cp.Mode)
	}
	if cp.StartPercent != 0.15 {
		t.Errorf("expected StartPercent=0.15, got %g", cp.StartPercent)
	}
	if cp.EndPercent != 0.95 {
		t.Errorf("expected EndPercent=0.95, got %g", cp.EndPercent)
	}
	if cp.ErrorDecayRate != 1 {
		t.Errorf("expected ErrorDecayRate=1, got %g", cp.ErrorDecayRate)
	}
	if !cp.UseRelativeThreshold {
		t.Errorf("expected UseRelativeThreshold=true, got %t", cp.UseRelativeThreshold)
	}
	if !cp.ResetErrorOnCompute {
		t.Errorf("expected ResetErrorOnCompute=true, got %t", cp.ResetErrorOnCompute)
	}
	if !cp.SCMPolicyDynamic {
		t.Errorf("expected SCMPolicyDynamic=true, got %t", cp.SCMPolicyDynamic)
	}
	if cp.SpectrumWarmupSteps != 4 {
		t.Errorf("expected SpectrumWarmupSteps=4, got %d", cp.SpectrumWarmupSteps)
	}
}

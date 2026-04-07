package gosd

import (
	"testing"
)

func TestGetSystemInfo(t *testing.T) {
	Load()
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

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

	return nil
}

func GetSystemInfo() string {
	var systemInfo *byte

	getSystemInfo.Call(unsafe.Pointer(&systemInfo))
	if systemInfo == nil {
		return ""
	}

	return charToString(systemInfo)
}

func Commit() string {
	var commitInfo *byte

	commit.Call(unsafe.Pointer(&commitInfo))
	if commitInfo == nil {
		return ""
	}

	return charToString(commitInfo)
}

func Version() string {
	var versionInfo *byte

	version.Call(unsafe.Pointer(&versionInfo))
	if versionInfo == nil {
		return ""
	}

	return charToString(versionInfo)
}

func GetNumPhysicalCores() int {
	var count int

	getNumPhysicalCores.Call(unsafe.Pointer(&count))
	return count
}

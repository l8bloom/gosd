package gosd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"unsafe"

	"github.com/jupiterrider/ffi"
)

func utilsGetNulString() *byte {
	return utilsStrToNulString("")
}

func utilsStrToNulString(text string) *byte {
	s := []byte(text + "\x00")
	return &s[0]
}

func charToString(text *byte) string {
	if text == nil {
		return ""
	}
	if *text == 0 {
		return ""
	}

	n := 0
	for ptr := unsafe.Pointer(text); *(*byte)(ptr) != 0; n++ {
		ptr = unsafe.Pointer(uintptr(ptr) + 1)
	}

	return string(unsafe.Slice(text, n))
}

func boolToByte(b bool) uint8 {
	if b {
		return uint8(1)
	}
	return uint8(0)
}

func byteToBool(b uint8) bool {
	return b > 0
}

func stringToChar(text string) *byte {
	if strings.IndexByte(text, 0) != -1 {
		panic(fmt.Errorf("stringToChar: %q already contains null byte", text))
	}
	a := make([]byte, len(text)+1)
	copy(a, []byte(text))

	return &a[0]
}

func loadError(name string, err error) error {
	fmt.Println(fmt.Errorf("could not load %q: %w", name, err).Error())
	return fmt.Errorf("could not load %q: %w", name, err)
}

func stableDiffusionLoadError() string {
	return `
Failed to load stablediffusion.so

It looks like dependent libraries (e.g. ggml.so) could not be found.

Fix options:
1. Put all .so/.dll files in the same directory and set:
	export LD_LIBRARY_PATH=$GOSD_DYN_LIB:$LD_LIBRARY_PATH (Linux)
	export PATH=$GOSD_DYN_LIB:$PATH (Windows)
	export DYLD_LIBRARY_PATH=$GOSD_DYN_LIB:$DYLD_LIBRARY_PATH (Mac)
`
}

func loadLibrary(lib string) (ffi.Lib, error) {
	path := os.Getenv("GOSD_DYN_LIB")
	if path == "" {
		return ffi.Lib{}, fmt.Errorf("gosd: %q env var undefined", "GOSD_DYN_LIB")
	}

	filename := getLibraryFilename(path, lib)
	if _, err := os.Stat(filename); err != nil {
		return ffi.Lib{}, fmt.Errorf("gosd: library not found at %q", filename)
	}

	l, err := ffi.Load(filename)
	if err != nil {
		err = errors.New(stableDiffusionLoadError())
	}
	return l, err
}

func getLibraryFilename(path, lib string) string {
	switch runtime.GOOS {
	case "linux", "freebsd":
		return filepath.Join(path, fmt.Sprintf("lib%s.so", lib))
	case "windows":
		return filepath.Join(path, fmt.Sprintf("%s.dll", lib))
	case "darwin":
		return filepath.Join(path, fmt.Sprintf("lib%s.dylib", lib))
	default:
		panic(fmt.Sprintf("OS %q not supported", runtime.GOOS))
	}
}

// Load loads the stable-diffusion.cpp shared library at runtime and all dependent libs
func Load() error {
	lib, err := loadLibrary("stable-diffusion")
	if err != nil {
		return err
	}
	if err := loadContextRoutines(lib); err != nil {
		return err
	}

	if err := loadImageRoutines(lib); err != nil {
		return err
	}

	if err := loadSystemRoutines(lib); err != nil {
		return err
	}

	if err := loadCallbacks(lib); err != nil {
		return err
	}

	if err := loadVideosRoutines(lib); err != nil {
		return err
	}

	return nil
}

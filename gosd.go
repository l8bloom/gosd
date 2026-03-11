package gosd

import (
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
	if b > 0 {
		return true
	}
	return false
}

func stringToChar(text string) *byte {
	if strings.IndexByte(text, 0) != -1 {
		panic(fmt.Errorf("stringToChar: %q already contains null byte\n", text))
	}
	a := make([]byte, len(text)+1)
	copy(a, []byte(text))

	return &a[0]
}

func loadError(name string, err error) error {
	fmt.Println(fmt.Errorf("could not load %q: %w", name, err).Error())
	return fmt.Errorf("could not load %q: %w", name, err)
}

func loadLibrary(lib string) (ffi.Lib, error) {
	path := os.Getenv("SD_DYN_LIB")
	if path != "" {
		path = os.Getenv("SD_DYN_LIB")
	}
	if path == "" {
		return ffi.Lib{}, fmt.Errorf("Can't find runtime stable-diffusion libraries")
	}

	filename := getLibraryFilename(path, lib)
	return ffi.Load(filename)
}

func getLibraryFilename(path, lib string) string {
	switch runtime.GOOS {
	case "linux", "freebsd":
		return filepath.Join(path, fmt.Sprintf("lib%s.so", lib))
	default:
		panic(fmt.Sprintf("OS %q not supported", runtime.GOOS))
	}
}

// Load loads the stable-diffusion.cpp shared library and all dependent libs
func Load() error {
	lib, err := loadLibrary("stable-diffusion")
	if err != nil {
		return loadError("stable-diffusion", err)
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

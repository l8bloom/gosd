package services

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/jupiterrider/ffi"
)

// creates C,Cpp-compatible empty string
func UtilsGetNulString() *byte {
	s := []byte("\x00")
	return &s[0]
}

func UtilsStrToNulString(text string) *byte {
	s := []byte(text + "\x00")
	return &s[0]
}

func utilsGetNulString() *byte {
	s := []byte("\x00")
	return &s[0]
}

func utilsStrToNulString(text string) *byte {
	s := []byte(text + "\x00")
	return &s[0]
}

func loadError(name string, err error) error {
	fmt.Println(fmt.Errorf("could not load %q: %w", name, err).Error())
	return fmt.Errorf("could not load %q: %w", name, err)
}

func loadLibrary(path, lib string) (ffi.Lib, error) {
	if path == "" && os.Getenv("SD_DYN_LIB") != "" {
		path = os.Getenv("SD_DYN_LIB")
	}
	if path == "" {
		return ffi.Lib{}, fmt.Errorf("Can't find runtime stable-diffusion libraries")
	}

	filename := getLibraryFilename(path, lib)
	return ffi.Load(filename)
}

// fetches the .so lib created (and named!) by the stable-diffusion cmake build
func getLibraryFilename(path, lib string) string {
	switch runtime.GOOS {
	case "linux", "freebsd":
		return filepath.Join(path, fmt.Sprintf("lib%s.so", lib))
	default:
		panic(fmt.Sprintf("OS %q not supported", runtime.GOOS))
	}
}

// Load loads the shared stable-diffusion.cpp libraries from the specified path.
func Load(path string) error {
	lib, err := loadLibrary(path, "stable-diffusion")
	if err != nil {
		return loadError("stable-diffusion", err)
	}
	// lib, err = loader.LoadLibrary(path, "ggml-base")
	// if err != nil {
	// 	return err
	// }

	// if err := loadGGMLBase(lib); err != nil {
	// 	return err
	// }

	// lib, err = loader.LoadLibrary(path, "llama")
	// if err != nil {
	// 	return err

	// }

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

	return nil
}

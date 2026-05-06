// example how to convert model tensor with gosd and stable-diffusion.cpp

package main

import (
	"os"

	sd "github.com/l8bloom/gosd/pkg/gosd"
)

func main() {
	// load dynamic libs of stable_diffusion.cpp and its deps
	if err := sd.Load(); err != nil {
		panic(err.Error())
	}

	modelPath := os.Getenv("MODEL_TO_CONVERT")
	vaePath := ""
	outputPath := "converted_model.gguf"
	outputType := sd.TypeQ2_K
	tensorTypeRules := ""
	convertName := false

	sd.Convert(modelPath, vaePath, outputPath, sd.SDType(outputType), tensorTypeRules, convertName)

}

// The most minimal example of how to generate an image.

package main

import (
	"fmt"
	"os"

	sd "github.com/l8bloom/gosd/pkg/gosd"
)

func main() {
	// load the dynamic libraries
	if err := sd.Load(); err != nil {
		panic(err.Error())
	}

	// create and conifgure the inference context
	ctxParams := sd.ContextParamsInit()
	ctxParams.DiffusionModelPath = os.Getenv("DIFFUSION_MODEL_PATH")
	ctxParams.VAEPath = os.Getenv("VAE_PATH")
	ctxParams.LLMPath = os.Getenv("LLM_PATH")

	ctx := sd.NewContext(ctxParams)
	defer sd.FreeCtx(ctx)

	imgParams := sd.ImageGenParamsInit()
	imgParams.Prompt = "a dog"

	genImage := sd.GenerateImage(ctx, imgParams)
	genImage.SavePNG("output.png")

	fmt.Printf("Image %q saved.\n", "output.png")
}

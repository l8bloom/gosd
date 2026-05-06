// example on how to use gosd for Canny processing
// 1st step: generate image
// 2nd step: apply canny algorithm

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
	ctxParams.DiffusionFlashAttn = true

	fmt.Printf("\nContext values:\n%s", sd.CtxParamsToStr(ctxParams))

	ctx := sd.NewContext(ctxParams)
	defer sd.FreeCtx(ctx)

	imgParams := sd.ImageGenParamsInit()

	// prompts
	imgParams.Prompt = "A simple indoor scene of a modern room with a chair, a table, and a window with sunlight coming through. Clean composition, strong lighting contrast, sharp edges, minimal clutter, photorealistic, centered composition, high detail."
	imgParams.NegativePrompt = "distortion, mascots, watermark, signature"

	imgParams.SampleParams.Guidance.TextCfg = 10
	imgParams.SampleParams.SampleSteps = 20
	imgParams.VAETilingParams.Enabled = true
	imgParams.VAETilingParams.RelSizeX = 4
	imgParams.VAETilingParams.RelSizeY = 4

	// image resolution
	imgParams.Width = 1536
	imgParams.Height = 768

	fmt.Printf("\nImage params:\n%s", sd.ImageGenParamsToStr(imgParams))

	genImage := sd.GenerateImage(ctx, imgParams)

	if err := genImage.SavePNG("output.png"); err != nil {
		fmt.Println("Saving image failed.")
		os.Exit(1)
	}
	fmt.Printf("Image %q saved.\n", "output.png")

	var highThreshold float32 = 0.05
	var lowThreshold float32 = 0.04
	var weak float32 = 0.8
	var strong float32 = 1.0
	inverse := false
	processed := sd.PreprocessCanny(
		genImage,
		highThreshold,
		lowThreshold,
		weak,
		strong,
		inverse,
	)

	if !processed {
		fmt.Println("Canny processing failed.")
		os.Exit(1)
	}
	if err := genImage.SavePNG("canny_output.png"); err != nil {
		fmt.Println("Saving image failed.")
		os.Exit(1)
	}

	fmt.Printf("Image %q saved.\n", "canny_output.png")
}

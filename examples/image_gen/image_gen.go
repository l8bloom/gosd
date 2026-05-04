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

	ctxParams.DiffusionFlashAttn = true // potential hardware optimizations
	// ctxParams.KeepClipOnCPU = true // in case of lower vram

	fmt.Printf("\nContext values:\n%s", sd.CtxParamsToStr(ctxParams))

	ctx := sd.NewContext(ctxParams)
	defer sd.FreeCtx(ctx)

	// fetch some default image values and further configure
	imgParams := sd.ImageGenParamsInit()

	// prompts
	imgParams.Prompt = "a beautiful landscape, ultra-detailed, 8K resolution, photorealistic, cinematic lighting"
	imgParams.NegativePrompt = "mascots, watermark, signature"
	// follow the instruction more closely
	imgParams.SampleParams.Guidance.TextCfg = 7

	// sampler config
	imgParams.SampleParams.SampleSteps = 10

	// vram saving configuration in case of lower vram
	imgParams.VAETilingParams.Enabled = true
	imgParams.VAETilingParams.RelSizeX = 4
	imgParams.VAETilingParams.RelSizeY = 4

	// image resolution
	imgParams.Width = 512
	imgParams.Height = 512

	// optionally refine/upscale the image after 1st generation pass
	// Hires = High Resolution
	imgParams.HiresParams.Enabled = true
	imgParams.HiresParams.Steps = 10
	// lower keeps it similar to 1st pass image, higher brings more variance
	imgParams.HiresParams.DenoisingStrength = 0.4
	imgParams.HiresParams.Scale = 2
	// HiresUpscalerLatent is the default mode
	// imgParams.HiresParams.Upscaler = sd.HiresUpscalerLatent

	fmt.Printf("\nImage params:\n%s", sd.ImageGenParamsToStr(imgParams))

	genImage := sd.GenerateImage(ctx, imgParams)
	genImage.SavePNG("output.png")
	fmt.Println(sd.Version())
	fmt.Println(sd.Commit())
}

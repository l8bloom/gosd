package main

import (
	"fmt"

	sd "github.com/l8bloom/gosd/pkg/gosd"
)

func main() {
	// load the dynamic libraries
	sd.Load()

	// create and conifgure the inference context
	ctxParams := sd.ContextParamsInit()

	ctxParams.DiffusionModelPath = "/tmp/stable.diffusion/flux-2-klein-9b-Q8_0.gguf"
	ctxParams.VAEPath = "/tmp/stable.diffusion/diffusion_pytorch_model.safetensors"
	ctxParams.LLMPath = "/tmp/stable.diffusion/Qwen3-8B-Q8_0.gguf"

	ctxParams.DiffusionFlashAttn = true // potential hardware optimizations
	// ctxParams.KeepClipOnCPU = true // in case of lower vram

	fmt.Printf("\nContext values:\n%s", sd.CtxParamsToStr(ctxParams))

	ctx := sd.NewContext(ctxParams)
	defer sd.FreeCtx(ctx)

	// fetch some default image values and further configure
	imgParams := sd.ImageGenParamsInit()

	// prompts
	imgParams.Prompt = "An orange cat on palm beach playing with oranges."
	imgParams.NegativePrompt = "mascots, watermark, signature"

	// sampler config
	imgParams.SampleParams.SampleSteps = 20

	// vram saving configuration in case of lower vram
	imgParams.VAETilingParams.Enabled = true
	imgParams.VAETilingParams.RelSizeX = 4
	imgParams.VAETilingParams.RelSizeY = 4

	// image resolution
	imgParams.Width = 1536
	imgParams.Height = 768

	fmt.Printf("\nImage params:\n%s", sd.ImageGenParamsToStr(imgParams))

	genImage := sd.GenerateImage(ctx, imgParams)
	genImage.SavePNG("output.png")
}

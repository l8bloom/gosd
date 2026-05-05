package main

import (
	"fmt"
	"os"

	sd "github.com/l8bloom/gosd/pkg/gosd"
)

// have a look at https://github.com/leejet/stable-diffusion.cpp/blob/master/README.md
// for list of supported models

func main() {
	// load the dynamic libraries
	if err := sd.Load(); err != nil {
		panic(err.Error())
	}

	// create and conifgure the inference context
	ctxParams := sd.ContextParamsInit()

	// https://huggingface.co/leejet/FLUX.2-klein-9B-GGUF/blob/main/flux-2-klein-9b-Q8_0.gguf
	ctxParams.DiffusionModelPath = os.Getenv("DIFFUSION_MODEL_PATH")

	// https://huggingface.co/black-forest-labs/FLUX.2-dev/blob/main/ae.safetensors
	ctxParams.VAEPath = os.Getenv("VAE_PATH")

	// https://huggingface.co/unsloth/Qwen3-8B-GGUF/blob/main/Qwen3-8B-Q8_0.gguf
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
	fmt.Printf("Image %q saved.\n", "output.png")
}

// example on how to use gosd to generate an image with stable-diffusion.cpp

package main

import (
	"fmt"
	"os"

	sd "github.com/l8bloom/gosd/pkg/gosd"
)

// have a look at the compatible models:
// https://github.com/leejet/stable-diffusion.cpp/blob/master/docs/wan.md

// NOTE: this script relies on ffmpeg to save the frames,
// without it present in the system PATH, generated video will be lost
func main() {
	// load dynamic libs of stable_diffusion.cpp and its deps
	if err := sd.Load(); err != nil {
		panic(err.Error())
	}

	// fetch default context values
	ctxParams := sd.ContextParamsInit()

	// https://huggingface.co/QuantStack/Wan2.2-T2V-A14B-GGUF/blob/main/LowNoise/Wan2.2-T2V-A14B-LowNoise-Q2_K.gguf
	ctxParams.DiffusionModelPath = os.Getenv("VIDEO_DIFFUSION_MODEL_PATH")

	// https://huggingface.co/QuantStack/Wan2.2-T2V-A14B-GGUF/blob/main/HighNoise/Wan2.2-T2V-A14B-HighNoise-Q2_K.gguf
	ctxParams.HighNoiseDiffusionModelPath = os.Getenv("VIDEO_HIGH_NOISE_DIFFUSION_MODEL_PATH")

	// https://huggingface.co/Comfy-Org/Wan_2.1_ComfyUI_repackaged/blob/main/split_files/vae/wan_2.1_vae.safetensors
	ctxParams.VAEPath = os.Getenv("VIDEO_VAE_PATH")

	// https://huggingface.co/city96/umt5-xxl-encoder-gguf/blob/main/umt5-xxl-encoder-Q8_0.gguf
	ctxParams.T5XXLPath = os.Getenv("VIDEO_T5XXL_PATH")

	ctxParams.DiffusionFlashAttn = true // enable potential hardware optimizations
	ctxParams.KeepClipOnCPU = true      // text encoding on cpu in case of lower vram

	// create context
	ctx := sd.NewContext(ctxParams)
	defer sd.FreeCtx(ctx)

	// grab some default video generation params
	vidParams := sd.VideoGenParamsInit()

	// optionally further customize the video params:

	// split spatial volume in case of lower vram
	vidParams.VAETilingParams.Enabled = true
	vidParams.VAETilingParams.RelSizeX = 4
	vidParams.VAETilingParams.RelSizeY = 4

	// low and high noise sampler params
	vidParams.SampleParams.SampleSteps = 10
	vidParams.SampleParams.SampleMethod = sd.EulerSampleMethod
	vidParams.SampleParams.Guidance.TextCfg = 6

	vidParams.HighNoiseSampleParams.SampleSteps = 10

	// number of video frames to generate
	vidParams.VideoFrames = 48

	// prompts
	vidParams.Prompt = "A narrow street in a rainy cyberpunk city at night. Neon signs reflect on wet pavement as light rain falls. A person with a transparent umbrella walks slowly past the camera while cars move in the distance. Steam rises from street vents and neon reflections ripple in puddles. Cinematic lighting, atmospheric fog, smooth camera pan, ultra detailed, realistic reflections."
	vidParams.NegativePrompt = "low quality, blurry, distorted, deformed, watermark, text, oversaturated, jpeg artifacts"

	// video resolution
	// vidParams.Width = 480
	// vidParams.Height = 832

	genVideo := sd.GenerateVideo(ctx, vidParams)
	framesPerSecond := 24
	genVideo.Save("output.mp4", framesPerSecond)
	fmt.Printf("Video %q saved, number of frames: %d\n", "output.mp4", len(genVideo.Data))
}

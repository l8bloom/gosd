package main

import (
	"fmt"

	sd "github.com/l8bloom/gosd/pkg/gosd"
)

func main() {
	// load the dynamic libs
	sd.Load()

	// fetch default context values
	ctxParams := sd.ContextParamsInit()

	// customize aspects of the context
	ctxParams.DiffusionModelPath = "/tmp/stable.diffusion/video/wan2.2_T2V_14B_q2_K/Wan2.2-T2V-A14B-LowNoise-Q2_K.gguf"
	ctxParams.HighNoiseDiffusionModelPath = "/tmp/stable.diffusion/video/wan2.2_T2V_14B_q2_K/Wan2.2-T2V-A14B-HighNoise-Q2_K.gguf"
	ctxParams.VAEPath = "/tmp/stable.diffusion/video/wan2.2_T2V_14B_q8/wan_2.1_vae.safetensors"
	ctxParams.LLMPath = "/tmp/stable.diffusion/video/wan2.2_T2V_14B_q8/umt5-xxl-encoder-Q8_0.gguf"
	ctxParams.T5XXLPath = "/tmp/stable.diffusion/video/wan2.2_T2V_14B_q8/umt5-xxl-encoder-Q8_0.gguf"
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
	vidParams.SampleParams.SampleSteps = 20
	vidParams.SampleParams.SampleMethod = sd.EulerSampleMethod
	vidParams.SampleParams.Guidance.TextCfg = 6

	vidParams.HighNoiseSampleParams.SampleSteps = 20

	// number of video frames to generate
	vidParams.VideoFrames = 48

	// prompts
	vidParams.Prompt = "A narrow street in a rainy cyberpunk city at night. Neon signs reflect on wet pavement as light rain falls. A person with a transparent umbrella walks slowly past the camera while cars move in the distance. Steam rises from street vents and neon reflections ripple in puddles. Cinematic lighting, atmospheric fog, smooth camera pan, ultra detailed, realistic reflections."
	vidParams.NegativePrompt = "low quality, blurry, distorted, deformed, watermark, text, oversaturated, jpeg artifacts"

	// video resolution
	vidParams.Width = 480
	vidParams.Height = 832

	genVideo := sd.GenerateVideo(ctx, vidParams)
	fmt.Println("Video generated", "Number of frames: ", len(genVideo.Data))
	framesPerSecond := 24
	genVideo.Save("output.mp4", framesPerSecond)
}

// example on how to use gosd to generate a video with audio
// NOTE: this doesn't work(atm) on ROCm/HIP backends

package main

import (
	"fmt"
	"os"
	"unsafe"

	sd "github.com/l8bloom/gosd/pkg/gosd"
)

var myLogCallback sd.LogCallback = func(level sd.LogLevel, text string, data unsafe.Pointer) {
	fmt.Println("My log callback:")
	fmt.Println("level: ", level)
	fmt.Println("text: ", text)
}

// NOTE: this script relies on ffmpeg to save the frames,
// without it present in the system PATH, generated video will be lost
func main() {
	// load dynamic libs of stable_diffusion.cpp and its deps
	if err := sd.Load(); err != nil {
		panic(err.Error())
	}
	sd.SetLogCallback(myLogCallback, nil)
	// fetch default context values
	ctxParams := sd.ContextParamsInit()

	// https://huggingface.co/unsloth/LTX-2.3-GGUF/blob/main/ltx-2.3-22b-dev-UD-Q4_K_M.gguf
	ctxParams.DiffusionModelPath = os.Getenv("VIDEO_EX_DIFFUSION_MODEL_PATH")

	// https://huggingface.co/unsloth/LTX-2.3-GGUF/blob/main/text_encoders/ltx-2.3-22b-dev_embeddings_connectors.safetensors
	ctxParams.EmbeddingsConnectorsPath = os.Getenv("VIDEO_EX_EMBEDDINGS_PATH")

	// https://huggingface.co/unsloth/gemma-3-12b-it-qat-GGUF/blob/main/gemma-3-12b-it-qat-UD-Q4_K_XL.gguf
	ctxParams.LLMPath = os.Getenv("VIDEO_EX_T5XXL_PATH")

	// https://huggingface.co/unsloth/LTX-2.3-GGUF/blob/main/vae/ltx-2.3-22b-dev_video_vae.safetensors
	ctxParams.VAEPath = os.Getenv("VIDEO_EX_VAE_PATH")

	// https://huggingface.co/unsloth/LTX-2.3-GGUF/blob/main/vae/ltx-2.3-22b-dev_audio_vae.safetensors
	ctxParams.AudioVAEPath = os.Getenv("AUDIO_EX_VAE_PATH")

	ctxParams.DiffusionFlashAttn = true // enable potential hardware optimizations

	// Optionally cherry-pick backend for each model/component
	// ctxParams.Backend = "all=vulkan0,vae=cpu"
	// ctxParams.ParamsBackend = "all=vulkan0,vae=cpu"

	// create context
	ctx := sd.NewContext(ctxParams)
	defer sd.FreeCtx(ctx)

	vidParams := sd.VideoGenParamsInit()

	// split spatial volume in case of lower vram
	vidParams.VAETilingParams.Enabled = true
	vidParams.VAETilingParams.RelSizeX = 64
	vidParams.VAETilingParams.RelSizeY = 64

	vidParams.SampleParams.SampleSteps = 30
	vidParams.SampleParams.SampleMethod = sd.EulerSampleMethod
	vidParams.SampleParams.Guidance.TextCfg = 6

	vidParams.FPS = 24

	// number of video frames to generate
	vidParams.VideoFrames = 240

	// prompts
	vidParams.Prompt = "A low-angle, fast-paced tracking shot following a massive stone golem awakening inside an ancient, crumbling moss-covered temple. Crimson magical energy cracks through the fissures of its rocky body as it slowly lifts its heavy head. The camera violently shakes as the golem slams its fist into the stone floor, sending debris and dust flying toward the lens. The audio features a loud, low-frequency stone-grinding rumble followed by a massive, echoing boom as the fist impacts the ground, accompanied by a tense orchestral string rise."
	vidParams.NegativePrompt = "low quality, blurry, distorted, deformed, watermark, text, oversaturated, jpeg artifacts"

	// video resolution
	vidParams.Width = 480
	vidParams.Height = 832

	// vidParams.Width = 240
	// vidParams.Height = 416

	sd.SetLogCallback(myLogCallback, nil)
	genVideo := sd.GenerateVideo(ctx, vidParams)
	framesPerSecond := 24
	genVideo.Save("output.mp4", framesPerSecond)
	fmt.Printf("Video %q saved, number of frames: %d\n", "output.mp4", len(genVideo.Data))

	sd.FreeAudio(&genVideo.Audio)
}

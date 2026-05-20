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
	vidParams.VAETilingParams.RelSizeX = 4
	vidParams.VAETilingParams.RelSizeY = 4

	vidParams.SampleParams.SampleSteps = 50
	vidParams.SampleParams.SampleMethod = sd.EulerSampleMethod
	vidParams.SampleParams.Guidance.TextCfg = 6

	vidParams.FPS = 24

	// number of video frames to generate
	vidParams.VideoFrames = 120

	// prompts
	vidParams.Prompt = "A cinematic, slow-motion shot of a narrow street in a rainy cyberpunk city at night. A person holding a transparent umbrella walks slowly past the camera. Neon signs reflect flawlessly on the wet pavement. Continuous light rain falls, creating ripples in puddles as steam rises from street vents and cars move in the far distance. Atmospheric fog, smooth camera pan, ultra-detailed realistic reflections. Concurrently, the synchronized audio track delivers the crisp, close-up acoustics of continuous soft rain drops falling, layered over a muffled, distant thunderstorm rumbling gently in the far background."
	vidParams.NegativePrompt = "low quality, blurry, distorted, deformed, watermark, text, oversaturated, jpeg artifacts"

	// video resolution
	vidParams.Width = 300
	vidParams.Height = 500

	sd.SetLogCallback(myLogCallback, nil)
	genVideo := sd.GenerateVideo(ctx, vidParams)
	framesPerSecond := 24
	genVideo.Save("output.mp4", framesPerSecond)
	fmt.Printf("Video %q saved, number of frames: %d\n", "output.mp4", len(genVideo.Data))

	sd.FreeAudio(&genVideo.Audio)
}

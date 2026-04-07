package main

import (
	"fmt"
	"unsafe"

	sd "github.com/l8bloom/gosd/pkg/gosd"
)

// all callbacks allow to be passed any type of app/user defined data
// but each callback must cast the data to its real type before using it
var myLogCallback sd.LogCallback = func(level sd.LogLevel, text string, data unsafe.Pointer) {
	fmt.Println("My log callback:")
	fmt.Println("level: ", level)
	fmt.Println("text: ", text)
}

var myProgressCallback sd.ProgressCallback = func(step int32, steps int32, time float32, data unsafe.Pointer) {
	fmt.Println("My progress callback:")
	fmt.Println("Step: ", step)
	fmt.Println("Steps: ", steps)
	fmt.Println("Time: ", time)
}

var myImagePreviewCallback sd.PreviewCallback[sd.Image] = func(step int32, image sd.Image, isNoisy bool, data unsafe.Pointer) {
	fmt.Println("My preview callback:")
	fmt.Println("Step: ", step)
	fmt.Println("Image size: ", image.Channel*image.Width*image.Height)
	fmt.Println("Still noisy?: ", isNoisy)

	(*ImageCtx)(data).ch <- image // cast data
	fmt.Printf("My preview callback sent to the channel, done for step %d\n", step)
}

func savePreviewImages(img *ImageCtx) {
	i := 0
	for {
		image := <-img.ch
		fmt.Println("Got the message from the channel")
		filename := fmt.Sprintf("output%03d.png", i)
		image.SavePNG(filename)
		fmt.Printf("\nImage %q saved.\n", filename)
		i++
	}
}

type ImageCtx struct {
	ch chan sd.Image
}

func main() {
	sd.Load()

	ctxParams := sd.ContextParamsInit()
	ctxParams.DiffusionModelPath = "/tmp/stable.diffusion/flux-2-klein-9b-Q8_0.gguf"
	ctxParams.VAEPath = "/tmp/stable.diffusion/diffusion_pytorch_model.safetensors"
	ctxParams.LLMPath = "/tmp/stable.diffusion/Qwen3-8B-Q8_0.gguf"

	ctxParams.DiffusionFlashAttn = true

	ctx := sd.NewContext(ctxParams)
	defer sd.FreeCtx(ctx)

	fmt.Printf("Context params to string:\n%s\n", sd.CtxParamsToStr(ctxParams))

	imgParams := sd.ImageGenParamsInit()

	imgParams.Prompt = "ultra detailed modern scene, a confused raccoon wearing a tiny business suit presenting a startup pitch on a large glass screen in a sleek futuristic office, holographic charts floating in the air, a robot intern taking notes on a tablet, a sleepy corgi wearing VR goggles under the desk, floor-to-ceiling windows showing a neon cyberpunk city at sunset, reflections on polished marble floor, coffee cups, sticky notes, laptop with glowing keyboard, whiteboard full of messy diagrams and equations, humorous tone, cinematic lighting, volumetric light rays, depth of field, sharp focus, intricate textures, photorealistic, 8k, highly detailed environment, modern tech aesthetic, global illumination, complex composition"
	imgParams.NegativePrompt = "blurry, low resolution, bad anatomy, extra limbs, distorted face, watermark, text artifacts, oversaturated, jpeg artifacts"

	imgParams.SampleParams.SampleSteps = 10

	imgParams.VAETilingParams.Enabled = true
	imgParams.VAETilingParams.RelSizeX = 4
	imgParams.VAETilingParams.RelSizeY = 4

	imgParams.Width = 1536
	imgParams.Height = 768

	data := &ImageCtx{ch: make(chan sd.Image)}

	// start the listener
	go savePreviewImages(data)

	// sd.SetLogCallback(myLogCallback, nil)
	// sd.SetProgressCallback(myProgressCallback, nil)
	sd.SetPreviewCallback(myImagePreviewCallback, sd.PreviewVAE, 1, true, false, unsafe.Pointer(data))

	genImage := sd.GenerateImage(ctx, imgParams)
	genImage.SavePNG("output.png")
}

// example on how callbacks can be used with gosd and stable-diffusion.cpp

package main

import (
	"fmt"
	"os"
	"unsafe"

	sd "github.com/l8bloom/gosd/pkg/gosd"
)

// all callbacks allow to be passed any type of data - note the unsafe Pointer(= void* in c/cpp)
// but each callback must cast the data to its real type, if it's using it
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

// deep-copy the image and send to a listener
var myImagePreviewCallback sd.PreviewCallback[sd.Image] = func(step int32, image sd.Image, isNoisy bool, data unsafe.Pointer) {
	fmt.Println("My preview callback:")
	fmt.Println("Step: ", step)
	fmt.Println("Image size: ", len(image.Data))
	fmt.Println("Still noisy?: ", isNoisy)

	(*ImageCtx)(data).ch <- image.Clone() // cast the data param!
	fmt.Printf("My preview callback sent to the channel, done for step %d\n", step)
}

func savePreviewImages(img *ImageCtx) {
	i := 0
	for {
		image := <-img.ch
		fmt.Println("\nGot the message from the channel")
		filename := fmt.Sprintf("output%03d.png", i+1)
		image.SavePNG(filename)
		fmt.Printf("\nImage %q saved.\n", filename)
		i++
	}
}

type ImageCtx struct {
	ch chan sd.Image
}

func main() {
	// load dynamic libs of stable_diffusion.cpp and its deps
	if err := sd.Load(); err != nil {
		panic(err.Error())
	}

	ctxParams := sd.ContextParamsInit()
	ctxParams.DiffusionModelPath = os.Getenv("DIFFUSION_MODEL_PATH")
	ctxParams.VAEPath = os.Getenv("VAE_PATH")
	ctxParams.LLMPath = os.Getenv("LLM_PATH")
	ctxParams.DiffusionFlashAttn = true

	ctx := sd.NewContext(ctxParams)
	defer sd.FreeCtx(ctx)

	fmt.Printf("Context params to string:\n%s\n", sd.CtxParamsToStr(ctxParams))

	imgParams := sd.ImageGenParamsInit()
	imgParams.Prompt = "A cinematic wide-angle view of a psychedelic cosmic garden; LEFT SIDE: towering giant translucent mushrooms glowing with neon cyan and violet iridescent fractals; CENTER: a sleek elongated tiger made of polished solid gold and glowing amber embers, leaping mid-air in a powerful stretch, clear muscular definition, high contrast; RIGHT SIDE: a shimmering ethereal fountain of liquid light and floating emerald lotus flowers; BACKGROUND: deep indigo space with sparkling star clusters, vibrant butterflies, hyper-detailed bioluminescent flora, whimsical atmosphere, masterpiece, 8k, joyful energy, sharp focus, ethereal lighting, magical realism."
	imgParams.NegativePrompt = "scary, monochromatic, blurry, low contrast, distorted, deformed, muddy colors, creepy, anatomical nonsense, watermark, grainy, jagged edges."

	imgParams.SampleParams.SampleSteps = 10
	imgParams.SampleParams.Guidance.TextCfg = 7

	imgParams.VAETilingParams.Enabled = true
	imgParams.VAETilingParams.RelSizeX = 4
	imgParams.VAETilingParams.RelSizeY = 4

	imgParams.Width = 1536
	imgParams.Height = 768

	data := &ImageCtx{ch: make(chan sd.Image)}

	// start the listener
	go savePreviewImages(data)

	// register callbacks
	// sd.SetLogCallback(myLogCallback, nil)
	// sd.SetProgressCallback(myProgressCallback, nil)
	sd.SetPreviewCallback(myImagePreviewCallback, sd.PreviewVAE, 1, true, false, unsafe.Pointer(data))

	genImage := sd.GenerateImage(ctx, imgParams)
	genImage.SavePNG("output.png")
}

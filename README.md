# gosd
![gosd](https://github.com/l8bloom/gosd/blob/main/assets/images/gosd.png)

Pure Go bindings for [stable-diffusion.cpp](https://github.com/leejet/stable-diffusion.cpp)

## Quick start

gosd allows writing Go programs backed by the stable-diffusion.cpp powerhouse.  
Getting started with it is mostly straightforward.

### Installation

```
go get github.com/l8bloom/gosd
```

To add the lib to project's dependency tree.

After that, the only thing left is to get the stable-diffusion shared libraries.  
There are multiple ways to do it, and is directed with the underlying OS and hardware.

The easiest way is to fetch an official release from the stable-diffuion project.


#### Manual installation

Here is a quick overview example of building the stack and sd libs on Linux with [Khronos Vulkan](https://www.vulkan.org/) gpu driver for Radeon.

```
# fetch the driver for your distribution, eg. on Ubuntu:
sudo apt update
sudo apt install mesa-vulkan-drivers vulkan-tools

# fetch the Vulkan loader
VULKAN_VER="1.4.341.1"
mkdir -p /tmp/vulkan && cd /tmp/vulkan \
    && wget -O /tmp/vulkan/vulkan.tar.xz "https://sdk.lunarg.com/sdk/download/${VULKAN_VER}/linux/vulkansdk-linux-x86_64-${VULKAN_VER}.tar.xz" \
    && tar -xf vulkan.tar.xz \
    && rm -rf vulkan.tar.xz

# update the env and run commands to confirm successful installation
. /tmp/vulkan/${VULKAN_VER}/setup-env.sh
vulkaninfo
vkcube

# clone the stable diffusion
git clone --recursive https://github.com/leejet/stable-diffusion.cpp
cd stable-diffusion.cpp

# build the sd
cmake -B builds/vulkan -DSD_VULKAN=ON -DSD_BUILD_SHARED_LIBS=ON -DSD_BUILD_SHARED_GGML_LIB=ON
cmake --build builds/vulkan --config Release

# gosd is ready now, export the lib root and try out examples
export GOSD_DYN_LIB="$(realpath builds/vulkan/bin/)"
git clone https://github.com/l8bloom/gosd && cd gosd
go run examples/image_gen/image_gen.go
```

## Examples


### Image generation

Example of image generation with gosd library.
```go
package main

import (
	"fmt"

	sd "github.com/l8bloom/gosd"
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
```

### Image generation with a preview

`examples/callbacks/image_gen/image_gen_with_callbacks.go` shows image generation with preview callback set.

Here are the `1st`, `4th` and last image generated:

![gosd](https://github.com/l8bloom/gosd/blob/main/assets/images/readmeFirstImage.png)
![gosd](https://github.com/l8bloom/gosd/blob/main/assets/images/readmeFourthImage.png)
![gosd](https://github.com/l8bloom/gosd/blob/main/assets/images/readmeFinalImage.png)

`examples` folder has more snippets showcasing classic use-cases.

[stable-diffusion.cpp](https://github.com/leejet/stable-diffusion.cpp) documentation
provides more insights into library's features, examples, models, licensing etc.


## Environment variables

- `GOSD_DYN_LIB` indicates root of stable-diffusion shared lib(.so, .dll etc.)

You may need to extend OS search path to load libraries sd depends on(eg on Linux `LD_LIBRARY_PATH`).


## Thanks

- C/Cpp: stable-diffusion.cpp, llama.cpp
- Go: purego, ffi


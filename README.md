# gosd
![gosd](https://github.com/l8bloom/gosd/blob/main/assets/images/gosd.png)

Pure Go bindings for [stable-diffusion.cpp](https://github.com/leejet/stable-diffusion.cpp)

## Quick start

gosd allows writing Go programs backed by the stable-diffusion.cpp powerhouse.  

### Installation

```
go get github.com/l8bloom/gosd
```

To add the lib to project's dependency tree.

After that, the only thing left is to get the stable-diffusion shared libraries.  
There are multiple ways to do it, and is directed with the underlying OS and hardware.

The easiest way is to fetch an official release from the stable-diffusion project.


#### Manual installation

Here is a quick overview example of building the stack and sd libs on Linux with [Khronos Vulkan](https://www.vulkan.org/) gpu driver for Radeon.

```bash
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
Getting started with gosd is rather straightforward.

### Image generation

```go
package main

import (
	"fmt"

	sd "github.com/l8bloom/gosd"
)

func main() {
	// load the dynamic libraries
	sd.Load()

	// create and configure the inference context
	ctxParams := sd.ContextParamsInit()

	// declare models
	ctxParams.DiffusionModelPath = "/tmp/stable.diffusion/flux-2-klein-9b-Q8_0.gguf"
	ctxParams.VAEPath = "/tmp/stable.diffusion/diffusion_pytorch_model.safetensors"
	ctxParams.LLMPath = "/tmp/stable.diffusion/Qwen3-8B-Q8_0.gguf"

	// create model's context
	ctx := sd.NewContext(ctxParams)
	defer sd.FreeCtx(ctx)

	// fetch some default image values and further configure
	imgParams := sd.ImageGenParamsInit()

	// prompts
	imgParams.Prompt = "An orange cat on palm beach playing with oranges."
	imgParams.NegativePrompt = "mascots, watermark, signature"

	genImage := sd.GenerateImage(ctx, imgParams)
	genImage.SavePNG("output.png")
}
```

Result:

![catImage](https://github.com/l8bloom/gosd/blob/main/assets/images/image_gen_ex_output.png)

### Image generation with a preview

`examples/callbacks/image_gen/image_gen_with_callbacks.go` shows image generation with a preview callback set.


#### 1st image
![image1](https://github.com/l8bloom/gosd/blob/main/assets/images/readmeFirstImage.png)
#### 4th image
![image4](https://github.com/l8bloom/gosd/blob/main/assets/images/readmeFourthImage.png)
#### 10th(last) image
![image10](https://github.com/l8bloom/gosd/blob/main/assets/images/readmeFinalImage.png)

`examples` folder has more snippets showcasing classic use-cases.

[stable-diffusion.cpp](https://github.com/leejet/stable-diffusion.cpp) documentation
provides more insights into library's features, examples, models etc.


## Environment variables

- `GOSD_DYN_LIB` indicates root of stable-diffusion shared lib(.so, .dll etc.)

You may need to extend OS search path to load libraries sd depends on(eg on Linux `LD_LIBRARY_PATH`).

## Portability

Built on top of purego and the inference libraries, gosd is portable across major systems  
and vast range of GPU/CPU hardware, but is being regularly tested only on linux platforms at the moment.


## Thanks

- C/Cpp: stable-diffusion.cpp, llama.cpp
- Go: ffi, purego
- Hugging Face

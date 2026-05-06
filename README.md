# gosd
![gosd](https://github.com/l8bloom/gosd/blob/main/assets/images/gosd.webp)

High-performance diffusion model inference in pure Go.

[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Go Reference](https://pkg.go.dev/badge/github.com/l8bloom/gosd.svg)](https://pkg.go.dev/github.com/l8bloom/gosd)
[![Linux](https://github.com/l8bloom/gosd/actions/workflows/linux.yaml/badge.svg)](https://github.com/l8bloom/gosd/actions/workflows/linux.yaml)
[![Windows](https://github.com/l8bloom/gosd/actions/workflows/windows.yaml/badge.svg)](https://github.com/l8bloom/gosd/actions/workflows/windows.yaml)
[![macOS](https://github.com/l8bloom/gosd/actions/workflows/macos.yaml/badge.svg)](https://github.com/l8bloom/gosd/actions/workflows/macos.yaml)
[![stable-diffusion.cpp](https://img.shields.io/badge/sd.cpp-6614334-yellow)](https://github.com/leejet/stable-diffusion.cpp/releases/tag/master-596-90e87bc)
[![Coverage](https://img.shields.io/badge/code%20coverage-80%25-purple)](https://github.com/l8bloom/gosd/actions)


## Features

- Image and video generation
- Image editing
- High-resolution upscaling (Neural ESRGAN models + Latent-space methods)
- Callback support for progressive previews during inference
- Model conversion (to SafeTensors / GGUF, optional VAE merging, tensor type rules)
- Hardware-accelerated inference (CUDA, Metal, Vulkan, ROCm and CPU)
- Minimal performance overhead compared to C/C++

## Quick start

gosd library is a set of pure Go bindings(no CGO) for [stable-diffusion.cpp](https://github.com/leejet/stable-diffusion.cpp).  
Equip your Go programs with image and video generation — hassle-free.

### Installation

```
go get github.com/l8bloom/gosd
```

to add gosd to your Go module.

After that, the only thing left is to get the stable-diffusion shared libraries. There are multiple ways to do it, and is directed with the underlying OS and hardware.

The simplest approach is to download an official release from the [stable-diffusion](https://github.com/leejet/stable-diffusion.cpp/releases) project that matches your system. The gosd library is designed to be agnostic regarding which specific build you choose though.  


<details>
<summary><strong>Manual installation (Linux + Vulkan example)</strong></summary>
<br>
Here is a quick overview example of building entire gpu stack, stable-diffusion libs and gosd on Linux with <a href="https://www.vulkan.org/">Khronos Vulkan</a> API for Radeon GPU.
<br>

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
</details>

## Examples

### Image generation

Generate an image in a few lines:

```go
package main

import sd "github.com/l8bloom/gosd/pkg/gosd"

func main() {
	// Load the dynamic libraries
	if err := sd.Load(); err != nil {
		panic(err.Error())
	}

	// Create and configure the inference context
	ctxParams := sd.ContextParamsInit()

	// Declare models
	ctxParams.DiffusionModelPath = "/tmp/stable.diffusion/flux-2-klein-9b-Q8_0.gguf"
	ctxParams.VAEPath = "/tmp/stable.diffusion/diffusion_pytorch_model.safetensors"
	ctxParams.LLMPath = "/tmp/stable.diffusion/Qwen3-8B-Q8_0.gguf"

	ctx := sd.NewContext(ctxParams)
	defer sd.FreeCtx(ctx)

	// Initialize image generation parameters
	imgParams := sd.ImageGenParamsInit()

	// Prompts
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

### ControlNet Preprocessing
Canny Edge Detection: Extract structural outlines from any image to guide ControlNet inference with pixel-perfect precision.

#### Generated Image
![cannyImage2](https://github.com/l8bloom/gosd/examples/system/canny/output.png)
#### Canny Preprocessed
![cannyImage2](https://github.com/l8bloom/gosd/examples/system/canny/canny_output.png)
#### ControlNet Output
![cannyImage3](https://github.com/l8bloom/gosd/examples/system/canny/image_from_canny_output.png)

(See `examples/system/canny/preprocess_canny.go` for a full implementation.)


`examples` folder has more snippets showcasing classic use-cases.

[stable-diffusion.cpp](https://github.com/leejet/stable-diffusion.cpp) documentation
provides more insights into library's features, examples, models etc.

## Environment variables

- `GOSD_DYN_LIB` indicates root of stable-diffusion shared lib(.so, .dll etc.)

You may need to extend OS search path to load libraries sd depends on.  
E.g. for Linux deployments: `export LD_LIBRARY_PATH=$GOSD_DYN_LIB:$LD_LIBRARY_PATH`

## Portability

CI/CD pipelines regularly test CPU-based inference on Linux, Windows, and macOS.  
GPU acceleration via Vulkan and AMD ROCm/HIP stack are being tested on privately hosted hardware, with verified compatibility on Linux.


| Platform | CPU (AMD64) | CPU (ARM64) | GPU (Vulkan 1.4.3) | GPU (ROCm 7.2.1) |
|----------|:-----------:|:-----------:|:------------------:|:----------------:|
| Linux    |      ✅     |      ✅     |         ✅         |         ✅       |
| Windows  |      ✅     |      ✅     |         -          |         -        |
| macOS    |      -      |      ✅     |         -          |         -        |

✅ = regularly tested

## Thanks

- C/C++: stable-diffusion.cpp, llama.cpp
- Go: ffi, purego
- Hugging Face

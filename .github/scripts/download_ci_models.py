import os
import time
from huggingface_hub import hf_hub_download
import requests

files_to_download = {
    "DIFFUSION_MODEL_PATH": (
        "unsloth/FLUX.2-klein-4B-GGUF",
        "flux-2-klein-4b-Q4_K_M.gguf",
    ),
    "VAE_PATH": ("ai-toolkit/flux2_vae", "ae.safetensors"),
    "LLM_PATH": ("unsloth/Qwen3-4B-GGUF", "Qwen3-4B-Q4_K_M.gguf"),
    "VIDEO_DIFFUSION_MODEL_PATH": (
        "QuantStack/Wan2.2-TI2V-5B-GGUF",
        "Wan2.2-TI2V-5B-Q2_K.gguf",
    ),
    "VIDEO_VAE_PATH": (
        "Comfy-Org/Wan_2.2_ComfyUI_Repackaged",
        "split_files/vae/wan2.2_vae.safetensors",
    ),
    "VIDEO_T5XXL_PATH": (
        "city96/umt5-xxl-encoder-gguf",
        "umt5-xxl-encoder-Q3_K_S.gguf",
    ),
    "MODEL_TO_CONVERT": (
        "black-forest-labs/FLUX.2-small-decoder",
        "diffusion_pytorch_model.safetensors",
    ),
}

# for non-huggingface sources
direct_downloads = {
    "UPSCALER_PATH": "https://github.com/xinntao/Real-ESRGAN/releases/download/v0.2.2.4/RealESRGAN_x4plus_anime_6B.pth"
}

results = {}

for env_var, (repo, filename) in files_to_download.items():
    print(f"Downloading {filename}...")

    start_time = time.time()
    absolute_path = hf_hub_download(repo_id=repo, filename=filename)
    duration = int(time.time() - start_time)

    results[env_var] = absolute_path
    print(f"Done in {duration}s | Stored at: {absolute_path}")

for env_var, url in direct_downloads.items():
    filename = url.split("/")[-1]
    save_path = os.path.abspath(filename)  # Saves to current CI directory

    print(f"Downloading {filename}...")

    start_time = time.time()
    with requests.get(url, stream=True) as r:
        r.raise_for_status()
        with open(save_path, "wb") as f:
            for chunk in r.iter_content(chunk_size=8192):
                f.write(chunk)

    duration = int(time.time() - start_time)
    print(f"Done in {duration}s | Stored at: {save_path}")
    results[env_var] = save_path

if "GITHUB_ENV" in os.environ:
    with open(os.environ["GITHUB_ENV"], "a") as f:
        for var, path in results.items():
            f.write(f"{var}={path}\n")

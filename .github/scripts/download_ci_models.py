import os
from huggingface_hub import hf_hub_download

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
        "vae/wan2.2_vae.safetensors",
    ),
    "VIDEO_T5XXL_PATH": (
        "city96/umt5-xxl-encoder-gguf",
        "umt5-xxl-encoder-Q3_K_S.gguf",
    ),
}

results = {}

for env_var, (repo, filename) in files_to_download.items():
    print(f"Downloading {filename}...")

    absolute_path = hf_hub_download(repo_id=repo, filename=filename)
    results[env_var] = absolute_path
    print(f"Stored at: {absolute_path}")

if "GITHUB_ENV" in os.environ:
    with open(os.environ["GITHUB_ENV"], "a") as f:
        for var, path in results.items():
            f.write(f"{var}={path}\n")

# %% [markdown]
# # Diffusers Model
#
# Stable Diffusion 系モデルの管理を担うオブジェクトについて定義します．
#
# from typing import Optional

# %%
import os
import re
import subprocess
import sys
from typing import Optional

import diffusers
import torch


# %%
class HuggingFaceDiffusionRepo:
    def __init__(
        self,
        reponame: str,
        revision: Optional[str] = None,
        torch_dtype=torch.float16,
        device: torch.device | str = "cuda",
    ):
        if revision is None:
            revision = self._get_latest_revision(reponame)
        self.reponame: str = reponame
        self.revision: str = revision
        self.torch_dtype = torch_dtype
        self.device = device

    @staticmethod
    def _get_latest_revision(reponame: str, host="https://huggingface.co") -> str:
        process = subprocess.Popen(
            [
                "git",
                "ls-remote",
                "--head",
                f"{host}/{reponame}.git",
            ],
            stdout=subprocess.PIPE,
            stderr=sys.stderr,
        )

        stdout, _ = process.communicate()
        return re.split(r"\t+", stdout.decode("ascii"))[0]

    def describe(self) -> dict[str, str]:
        return {
            "reponame": self.reponame,
            "revision": self.revision,
            "torch_dtype": str(self.torch_dtype),
            "device": str(self.device),
            "diffusers.__version__": diffusers.__version__,
        }

    def StableDiffusionPipeline(self):
        return diffusers.StableDiffusionPipeline.from_pretrained(
            self.reponame,
            revision=self.revision,
            torch_dtype=self.torch_dtype,
            cache_dir=os.getenv("HF_DATASETS_CACHE"),
        ).to(self.device)


# %%
if __name__ == "__main__":
    import json

    repo = HuggingFaceDiffusionRepo("852wa/8528-diffusion", torch_dtype=torch.float32)
    print(json.dumps(repo.describe()))

    with torch.autocast("cuda"):
        result = repo.StableDiffusionPipeline()(["pipe girl"])

    if result["nsfw_content_detected"][0] == True:
        with torch.autocast("cuda"):
            result = repo.StableDiffusionPipeline()(["pipe girl"])

    print(result)
    display(result["images"][0])


# %%
def execute_txt2img(
    model,
    prompt: str,
    negative_prompt: str,
    guidance_scale: float = 7.5,
    num_inference_steps: int = 50,
    width: int = 512,
    height: int = 512,
    seed: int = -1,
):
    if seed < 0:
        seed = int(torch.randint(0, 2**62, (1,)))
    generator = torch.Generator(device=model.device).manual_seed(seed)

    with torch.autocast("cuda"):
        result = model(
            prompt=[prompt],
            negative_prompt=[negative_prompt],
            guidance_scale=guidance_scale,
            num_inference_steps=num_inference_steps,
            generator=generator,
            width=width,
            height=height,
        )

    if "nsfw_content_detected" in result:
        if result["nsfw_content_detected"][0] == True:
            raise RuntimeError(f"detected NSFW content: {result}")

    return result["images"][0], {
        "prompt": prompt,
        "negative_prompt": negative_prompt,
        "seed": seed,
        "guidance_scale": guidance_scale,
        "num_inference_steps": num_inference_steps,
        "width": width,
        "height": height,
    }


# %%
if __name__ == "__main__":
    model = repo.StableDiffusionPipeline()
    model.safety_checker = None
    image, props = execute_txt2img(
        model,
        "happy valentine's day",
        "worse quality",
    )
    props["model"] = repo.describe()
    print(json.dumps(props))
    display(image)

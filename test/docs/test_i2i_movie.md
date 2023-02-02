---
jupyter:
  jupytext:
    cell_metadata_filter: -all
    formats: '@/ipynb,docs//md,py:percent'
    main_language: python
    text_representation:
      extension: .md
      format_name: markdown
      format_version: '1.3'
      jupytext_version: 1.14.4
---

```python
import diffusers
import torch
import PIL
from PIL import Image
import os
import numpy as np
from typing import List, Optional
```

```python
modelname = ""
prompt = ""
n_prompt = ""
width = 640
height = 480
inference_steps = 50
i2i_inference_steps = 30
strength = 0.3

def t2i():
    global modelname, prompt, n_prompt, width, height, inference_steps, i2i_inference_steps, strength

    t2i = diffusers.StableDiffusionPipeline.from_pretrained(
        modelname, cache_dir=os.getenv("HF_DATASETS_CACHE")).to("cuda")
    t2i.safety_checker = None

    try:
        with torch.autocast("cuda"):
            torch.cuda.empty_cache()
            return t2i(
                prompt=[prompt],
                negative_prompt=[n_prompt],
                num_inference_steps=inference_steps,
                width=640,
                height=480,
                generator=torch.Generator(device=t2i.device).manual_seed(2^63),
            )["images"][0]
    finally:
        del t2i
        torch.cuda.empty_cache()


img = t2i()
img
```

```python
def pil2cv(image):
    new_image = np.array(image, dtype=np.uint8)
    if new_image.ndim == 2:  # モノクロ
        pass
    elif new_image.shape[2] == 3:  # カラー
        new_image = cv2.cvtColor(new_image, cv2.COLOR_RGB2BGR)
    elif new_image.shape[2] == 4:  # 透過
        new_image = cv2.cvtColor(new_image, cv2.COLOR_RGBA2BGRA)
    return new_image

import cv2
from tqdm import tqdm
def i2i():
    global modelname, prompt, n_prompt, width, height, inference_steps, i2i_inference_steps, img
    i2i = diffusers.StableDiffusionImg2ImgPipeline.from_pretrained(
        modelname, cache_dir=os.getenv("HF_DATASETS_CACHE")).to("cuda")
    i2i.safety_checker = None
    i2i._progress_bar_config = {
        "disable": True,
    }

    writer = cv2.VideoWriter("video.mp4", cv2.VideoWriter_fourcc('m', 'p', '4', 'v'), 8,
                            (width, height))
    try:
        writer.write(pil2cv(img))

        for _ in tqdm(range(8 * 30)):
            # img = img.rotate(3, resample=Image.Resampling.BICUBIC)
            # clip_rate = 0.07
            # clip_horizonal = int(img.width * clip_rate * 0.5)
            # clip_vertical = int(img.height * clip_rate * 0.5)
            # img = img.crop((clip_horizonal, clip_vertical, img.width - clip_horizonal,
            #                 img.height - clip_vertical))
            # img = img.resize((width, height))
            
            # [0.1, 0.6) の範囲で線形降下する確率密度関数
            strength = 0.3 + (1 - np.sqrt(np.random.random())) * 0.2

            torch.cuda.empty_cache()
            with torch.autocast("cuda"):
                img = i2i(
                    prompt=[prompt],
                    negative_prompt=[n_prompt],
                    num_inference_steps=int(i2i_inference_steps / strength)+1,
                    image=img.convert("RGB"),
                    strength=strength,
                )["images"][0]
                writer.write(pil2cv(img))
    finally:
        writer.release()
        del i2i
        torch.cuda.empty_cache()

i2i()
```

---
jupyter:
  jupytext:
    formats: '@/ipynb,docs//md,py:percent'
    text_representation:
      extension: .md
      format_name: markdown
      format_version: '1.3'
      jupytext_version: 1.14.1
  kernelspec:
    display_name: .container-venv
    language: python
    name: python3
---

```python jupyter={"source_hidden": true} tags=[]
import diffusers
import torch
import os
from typing import List, Optional
from ipywidgets import widgets, interact, Layout
from PIL import Image

modelnames = [
    "stabilityai/stable-diffusion-2",
    "alfredplpl/cool-japan-diffusion-for-learning-2-0",
    "Linaqruf/anything-v3.0",
    "852wa/8528-diffusion",
    "hakurei/waifu-diffusion",
    "../models/8528d-v5",
    "../models/8528d-v4",
    "../models/8528d-v3",
    "../models/8528d-v2",
]

def gen_txt2img(pipe, prompt: str, negative_prompt: Optional[str] = None, width=512, height=512):
    with torch.autocast("cuda"):
        result = pipe([prompt],
                            negative_prompt=[negative_prompt] if negative_prompt is not None else None,
                            num_inference_steps=50,
                            width=width,
                            height=height)
        return result["images"][0]

def grid_img(imgs, rows: int, cols: int):
    if len(imgs) > rows * cols:
        raise Exception("rows * cols < len(imgs)")

    w, h = imgs[0].size
    grid = Image.new("RGB", size=(cols * w, rows * h))

    for i, img in enumerate(imgs):
        grid.paste(img, box=(i % cols * w, i // cols * h))

    return grid


n = 4
pipe = None
imgs = []
pipemodel = ""

```

```python tags=[]
import sys
import matplotlib.pyplot as plt
from IPython.display import clear_output, display

text_layout = Layout(width="80%")
slider_layout = Layout(width="60%")
model = widgets.Select(options=modelnames,
                       description="model",
                       ensure_option=True,
                       layout=slider_layout)
prompt = widgets.Text(placeholder="positive prompt here", layout=text_layout)
negative_prompt = widgets.Text(placeholder="negative prompt here",
                               layout=text_layout)
width = widgets.IntSlider(value=512,
                          min=360,
                          max=640,
                          step=2,
                          description="width",
                          layout=slider_layout)
height = widgets.IntSlider(value=512,
                           min=360,
                           max=640,
                           step=2,
                           description="height",
                           layout=slider_layout)
run = widgets.Button(description="run")

def display_widgets():
    global prompt, negative_prompt, width, height, run
    display(model)
    display(prompt)
    display(negative_prompt)
    display(width)
    display(height)
    display(run)


def set_model(clicked):
    global pipe, model
    if str(model.value) in set(modelnames):
        clear_output()
        display_widgets()
        pipe = diffusers.StableDiffusionPipeline.from_pretrained(
            str(model.value),
            # torch_dtype=torch.float16,
            cache_dir=os.getenv("HF_DATASETS_CACHE")).to("cuda")
        pipe.safety_checker = None

def generate(clicked):
    global prompt, negative_prompt, width, height, finished, pipe, pipemodel
    if pipe is None or str(model.value) != pipemodel:
        set_model(None)
        pipemodel = str(model.value)
    clear_output()
    display_widgets()
    print("\n".join([
        "```md",
        "# prompt",
        f"{str(prompt.value)}",
        "",
        "# negative prompt",
        f"{str(negative_prompt.value)}",
        "```",
    ]))

    imgs = []
    for i in range(4):
        torch.cuda.empty_cache()
        imgs.append(gen_txt2img(pipe,
                    str(prompt.value),
                    negative_prompt=str(negative_prompt.value),
                    width=width.value,
                    height=height.value))
        display(imgs[i])
    torch.cuda.empty_cache()

    display(grid_img(imgs, 2, 2))

run.on_click(generate)
display_widgets()
```

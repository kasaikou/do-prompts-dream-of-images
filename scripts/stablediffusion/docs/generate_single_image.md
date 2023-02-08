# シングル画像生成 `HuggingFace & diffusers` 編

```python
import json
from typing import Optional

import IPython.display
import ipywidgets.widgets
import PIL.ImageEnhance
import stablediffusion
import torch
```

```python
class Application:
    def __init__(self) -> None:
        from ipywidgets.widgets import (
            Button,
            FloatSlider,
            HBox,
            IntSlider,
            IntText,
            Layout,
            Text,
            VBox,
        )

        self.modelInfo: Optional[stablediffusion.HuggingFaceDiffusionRepo] = None
        self.model = None

        self.reponame = Text(description="reponame", placeholder="username/reponame")
        self.revision = Text(
            description="reponame",
            placeholder="if this is empty, select latest version",
        )
        self.prompt = Text(
            description="prompt",
            placeholder="happy valentine's day",
            layout=Layout(width="80%"),
        )
        self.negativePrompt = Text(
            description="n_prompt",
            placeholder="worse quality",
            layout=Layout(width="80%"),
        )
        self.width = IntSlider(description="width", value=512, min=360, max=640, step=2)
        self.height = IntSlider(
            description="height", value=512, min=360, max=640, step=2
        )
        self.inferenceSteps = IntSlider(
            description="inference steps", value=50, min=1, max=200, step=1
        )
        self.saturation = FloatSlider(
            description="saturation", value=1.0, min=0.0, max=2.0
        )
        self.contrast = FloatSlider(description="contrast", value=1.0, min=0.0, max=2.0)
        self.brightness = FloatSlider(
            description="brightness", value=1.0, min=0.0, max=2.0
        )

        self.seed = Text(description="seed", value="-1", allow_none=False)
        self.run = Button(description="run")
        self.run.on_click(self.execute)

        self.controller = VBox(
            [
                HBox([self.reponame, self.revision]),
                self.prompt,
                self.negativePrompt,
                HBox([self.width, self.height, self.inferenceSteps]),
                HBox([self.saturation, self.contrast, self.brightness]),
                HBox([self.seed, self.run]),
            ],
        )

        IPython.display.display(self.controller)

    def execute(self, clicked):
        print("load model information")

        selectModel = stablediffusion.HuggingFaceDiffusionRepo(
            reponame=str(self.reponame.value),
            revision=str(self.revision.value) if self.revision.value != "" else None,
            torch_dtype=torch.float32,
        )

        if self.modelInfo:
            updateModel = (
                selectModel.reponame != self.modelInfo.reponame
                or selectModel.revision != self.modelInfo.revision
            )
        else:
            updateModel = True

        if self.model is not None and updateModel:
            del self.model
            self.model = None
            torch.cuda.empty_cache()

        if self.model is None or updateModel:
            IPython.display.clear_output(False)
            IPython.display.display(self.controller)
            print("load model")
            self.modelInfo = selectModel
            self.model = self.modelInfo.StableDiffusionPipeline()
            self.model.safety_checker = None
        else:
            print("use model from cache")

        if self.modelInfo is None:
            raise RuntimeError("self.modelInfo is None")

        IPython.display.clear_output(False)
        IPython.display.display(self.controller)
        seed = int(str(self.seed.value)) if str(self.seed.value) != "" else -1

        print(json.dumps(self.modelInfo.describe()))
        for i in range(4):
            image, props = stablediffusion.execute_txt2img(
                model=self.model,
                prompt=str(self.prompt.value),
                negative_prompt=str(self.negativePrompt.value),
                num_inference_steps=self.inferenceSteps.value,
                width=self.width.value,
                height=self.height.value,
                seed=seed,
            )
            image = PIL.ImageEnhance.Color(image).enhance(self.saturation.value)
            props["saturation"] = self.saturation.value
            image = PIL.ImageEnhance.Contrast(image).enhance(self.contrast.value)
            props["contrast"] = self.contrast.value
            image = PIL.ImageEnhance.Brightness(image).enhance(self.brightness.value)
            props["brightness"] = self.brightness.value

            print(json.dumps(props))
            IPython.display.display(image)
            if seed > -1:
                break


app = Application()
```

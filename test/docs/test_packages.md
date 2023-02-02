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

```

```python
import torch
```

```python
# torch (cuda)
print(f"torch status: {torch.__version__=}")
print(f"{torch.cuda.is_available()=}")
if torch.cuda.is_available():
    print(f"{torch.cuda.get_device_name()=}")
    print(f"{torch.cuda.get_device_capability()=}")
```

```python
from ipywidgets import widgets

button = widgets.Button(description = "ok?")
display(button)
```

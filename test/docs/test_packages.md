---
jupyter:
  jupytext:
    cell_metadata_filter: -all
    formats: '@/ipynb,docs//md,py:percent'
    text_representation:
      extension: .md
      format_name: markdown
      format_version: '1.3'
      jupytext_version: 1.14.1
  kernelspec:
    display_name: 'Python 3.10.6 (''.container-venv'': venv)'
    language: python
    name: python3
---

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

# ---
# jupyter:
#   jupytext:
#     cell_metadata_filter: -all
#     formats: '@/ipynb,docs//md,py:percent'
#     text_representation:
#       extension: .py
#       format_name: percent
#       format_version: '1.3'
#       jupytext_version: 1.14.4
#   kernelspec:
#     display_name: .venv
#     language: python
#     name: python3
# ---

# %%
import torch
from ipywidgets import widgets

# %%
# torch (cuda)
print(f"torch status: {torch.__version__=}")
print(f"{torch.cuda.is_available()=}")
if torch.cuda.is_available():
    print(f"{torch.cuda.get_device_name()=}")
    print(f"{torch.cuda.get_device_capability()=}")

# %%

button = widgets.Button(description="ok?")
display(button)

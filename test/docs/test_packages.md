```python
import torch
from ipywidgets import widgets
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

button = widgets.Button(description="ok?")
display(button)
```

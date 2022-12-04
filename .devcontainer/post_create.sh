#!/bin/sh

cd /workspace &&
python3.10 -m venv .container-venv &&
.container-venv/bin/pip install -r requirements.lock --extra-index-url https://download.pytorch.org/whl/cu113

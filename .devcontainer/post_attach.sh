#!/bin/sh

cd /workspace
poetry env use python3.10

exec-onchanges \
    -i "*.py" -i "**/@*.ipynb" \
    -e .git -e __pycache__ -e .venv -e node_modules -e models -e repos -e .cache -e .ipynb_checkpoints -- \
    go run ./cmd/update_notebook "{{FILEPATH}}"
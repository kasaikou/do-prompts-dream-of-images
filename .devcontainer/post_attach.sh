#!/bin/sh

cd /workspace
exec-onchanges \
    -i "*.py" -i "**/@*.ipynb" \
    -e .git -e __pycache__ -e .container_venv -e node_modules -e models -e repos -e .cache -e .ipynb_checkpoints -- \
    .container-venv/bin/jupytext --set-formats @/ipynb,docs//md:markdown,py:percent "{{FILEPATH}}"

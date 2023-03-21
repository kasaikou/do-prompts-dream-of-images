#!/bin/sh

cd /workspace
poetry env use python3.10

exec-onchanges -i **.ipynb -e mnt -e .venv -e .git -- sh .devcontainer/format_notebook.sh "{{FILEPATH}}"
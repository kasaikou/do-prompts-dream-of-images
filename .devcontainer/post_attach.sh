#!/bin/sh

cd /workspace
poetry env use python3.10

exec-onchanges \
    -i "*.py" -i "**/@*.ipynb" \
    -e .git -e __pycache__ -e .venv -e node_modules -e models -e repos -e .cache -e .ipynb_checkpoints -- \
    .venv/bin/jupytext \
        --opt cell_metadata_filter=-all \
        --opt notebook_metadata_filter=-all,-jupytext.text_representation \
        --set-formats @/ipynb,docs//md:markdown,py:percent "{{FILEPATH}}"

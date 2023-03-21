#/bin/sh

FILENAME=$1
/workspace/.venv/bin/nbqa isort ${FILENAME} --float-to-top
/workspace/.venv/bin/nbqa black ${FILENAME}
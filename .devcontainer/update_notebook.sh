#!/bin/bash

cd /workspace

filepath=$1
filename=${filepath##*/}
extname=${filepath##*.}

if [ "${extname}" == "ipynb" ]; then
    pyfilename=`echo $filename | sed -e s%\.[^.]*$%.py%g | sed -e s%^@%%g`
    pyfilepath=`echo $filepath | sed -e s%${filename}%${pyfilename}%g`
    
elif [ "${extname}" == "py" ]; then
    pyfilename=filename
    pyfilepath=filepath

else
    echo "unknown ext: $extname"
    exit 1
fi

.venv/bin/jupytext \
    --pipe '.venv/bin/isort - --treat-comment-as-code "# %%" --float-to-top' \
    --pipe '.venv/bin/black -' \
    --opt cell_metadata_filter=-all \
    --opt notebook_metadata_filter=-all,-jupytext.text_representation \
    --set-formats @/ipynb,docs//md:markdown,py:percent $filepath

if [ "${filepath}" != "${pyfilepath}" ]; then
    .venv/bin/jupytext \
    --pipe '.venv/bin/isort - --treat-comment-as-code "# %%" --float-to-top' \
    --pipe '.venv/bin/black -' \
    --opt cell_metadata_filter=-all \
    --opt notebook_metadata_filter=-all,-jupytext.text_representation \
    --set-formats @/ipynb,docs//md:markdown,py:percent $pyfilepath
    
fi

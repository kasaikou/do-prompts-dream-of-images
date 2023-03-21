#!/bin/sh

go install github.com/streamwest-1629/exec-onchanges/cmd/exec-onchanges@latest

poetry install --only editor &&
/workspace/.venv/bin/nbstripout --install

git config --global filter.nbstripout.extrakeys '
  metadata.celltoolbar
  metadata.kernelspec
  metadata.language_info.codemirror_mode.version
  metadata.language_info.pygments_lexer
  metadata.language_info.version
  metadata.toc
  metadata.notify_time
  metadata.varInspector
  cell.metadata.heading_collapsed
  cell.metadata.hidden
  cell.metadata.code_folding
  cell.metadata.tags
  cell.metadata.init_cell'

#!/bin/sh

apt-get install -y \
    python3.10 \
    python3-pip \
    python3-venv \
    nodejs \
    npm \
    zlib1g-dev \
    libgif-dev \
    libjpeg-dev \
    libcairo2-dev \
    graphviz &&

wget -qO- "https://go.dev/dl/go1.19.linux-$(dpkg --print-architecture).tar.gz" | tar zxf - -C /usr/local
git lfs install --skip-smudge

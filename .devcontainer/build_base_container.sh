#!/bin/sh

# Install apt packages
apt-get update &&
apt-get install -y \
    wget \
    curl \
    tar \
    unzip \
    git \
    git-lfs \
    software-properties-common \
    build-essential \
    pkg-config

# # Download and Install cuda and cuda toolkit
# apt-key del 7fa2af80 &&
# wget https://developer.download.nvidia.com/compute/cuda/repos/wsl-ubuntu/x86_64/cuda-wsl-ubuntu.pin &&
# mv cuda-wsl-ubuntu.pin /etc/apt/preferences.d/cuda-repository-pin-600 &&
# wget https://developer.download.nvidia.com/compute/cuda/11.8.0/local_installers/cuda-repo-wsl-ubuntu-11-8-local_11.8.0-1_amd64.deb &&
# dpkg -i cuda-repo-wsl-ubuntu-11-8-local_11.8.0-1_amd64.deb &&
# cp /var/cuda-repo-wsl-ubuntu-11-8-local/cuda-*-keyring.gpg /usr/share/keyrings/ &&
# apt-get update &&
# apt-get -y install cuda cuda-toolkit-11-8

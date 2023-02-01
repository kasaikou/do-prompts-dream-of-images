ARG username=vscode
ARG useruid=1000
ARG usergid=${useruid}

FROM nvidia/cuda:11.8.0-cudnn8-devel-ubuntu22.04 as core

ARG username
ARG useruid
ARG usergid

ENV DEBIAN_FRONTEND=nointeractive
ENV TZ=Asia/Tokyo

# Create non-root user
RUN sed -i 's@archive.ubuntu.com@ftp.jaist.ac.jp/pub/Linux@g' /etc/apt/sources.list && \
    apt-get update
RUN apt-get install -y --no-install-recommends \
    sudo \
    tzdata \
    wget \
    curl \
    tar \
    unzip \
    software-properties-common && \
    groupadd --gid ${usergid} ${username} && \
    useradd -s /bin/bash --uid ${useruid} --gid ${usergid} -m ${username} && \
    echo ${username} ALL=\(root\) NOPASSWD:ALL > /etc/sudoers.d/${username} && \
    chmod 0440 /etc/sudoers.d/${username}    


FROM core as published

RUN apt-get install -y --no-install-recommends \
    git \
    git-lfs \
    imagemagick \
    build-essential \
    pkg-config && \
    ln -s /usr/bin/python3 /usr/bin/python

ARG username
USER ${username}
RUN curl -sSL https://install.python-poetry.org | python3 -

ENV SHELL=/bin/bash

ARG username=vscode
ARG useruid=1000
ARG usergid=${useruid}
ARG golang=1.19.5
ARG node=18.13.0

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


FROM core AS downloader

WORKDIR /opt
RUN mkdir -p /opt/usr/local && \
    chmod -R 755 /opt/usr/local

ARG golang
ARG node
RUN \
    # golang
    wget -q https://go.dev/dl/go${golang}.linux-amd64.tar.gz && \
    tar -C /opt/usr/local -xzf go${golang}.linux-amd64.tar.gz && \
    # node
    wget https://nodejs.org/download/release/v${node}/node-v${node}-linux-x64.tar.gz && \
    tar -C /opt/usr/local -xvf node-v${node}-linux-x64.tar.gz && \
    mv /opt/usr/local/node-v${node}-linux-x64/ /opt/usr/local/node/

FROM core as published

RUN apt-get install -y --no-install-recommends \
    git \
    git-lfs \
    python3.10 \
    imagemagick \
    build-essential \
    pkg-config && \
    update-alternatives --install /usr/bin/python python /usr/bin/python3.10 10 && \
    update-alternatives --install /usr/bin/python3 python3 /usr/bin/python3.10 10

ARG username
USER ${username}
RUN curl -sSL https://install.python-poetry.org | python3 -

COPY --from=downloader /opt/usr/local /usr/local

ENV GOPATH=/home/${username}/go
ENV PATH=${PATH}:/usr/local/go/bin
ENV PATH=${PATH}:${GOPATH}/bin
ENV PATH=${PATH}:/usr/local/node/bin

RUN \
    # Check version
    python --version && \
    node --version && \
    go version && \
    git --version

ENV SHELL=/bin/bash

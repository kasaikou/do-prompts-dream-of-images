ARG username=vscode
ARG useruid=1000
ARG usergid=${useruid}

FROM ubuntu:22.04 as based

ARG username
ARG useruid
ARG usergid

# Create non-root user
RUN apt-get update && \
    apt-get install -y sudo && \
    groupadd --gid ${usergid} ${username} && \
    useradd -s /bin/bash --uid ${useruid} --gid ${usergid} -m ${username} && \
    echo ${username} ALL=\(root\) NOPASSWD:ALL > /etc/sudoers.d/${username} && \
    chmod 0440 /etc/sudoers.d/${username}    

# Build container with shellscript
ARG scriptfile=build_base_container.sh
COPY .devcontainer/${scriptfile} .
USER ${username}
RUN sudo sed -i 's@archive.ubuntu.com@ftp.jaist.ac.jp/pub/Linux@g' /etc/apt/sources.list && \
    sudo chown ${username} ./${scriptfile} && \
    sudo chmod +x ./${scriptfile} && \
    sudo bash ./${scriptfile}

FROM based as app

ARG username

ARG scriptfile=build_app_container.sh
COPY .devcontainer/${scriptfile} .
USER ${username}
RUN sudo chown ${username} ./${scriptfile} && \
    sudo chmod +x ./${scriptfile} && \
    sudo bash ./${scriptfile}

ENV DEBIAN_FRONTEND=nointeractive
ENV LANG=C.UTF-8
ENV SHELL=/bin/bash

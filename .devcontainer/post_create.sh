#!/bin/sh

cd /workspace &&
python3.10 -m venv .container-venv &&
.container-venv/bin/pip install -r requirements.lock --extra-index-url https://download.pytorch.org/whl/cu113

# install go development kit
cd /workspace/tools
sudo chown -R vscode ~/go
go install golang.org/x/tools/gopls@latest
go install golang.org/x/lint/golint@latest
go install github.com/go-delve/delve/cmd/dlv@master
go install github.com/haya14busa/goplay/cmd/goplay@v1.0.0
go install github.com/fatih/gomodifytags@v1.16.0
go install github.com/josharian/impl@latest
go install github.com/cweill/gotests/gotests@latest
go install github.com/ramya-rao-a/go-outline@latest
go install golang.org/x/tools/cmd/godoc@latest
go install honnef.co/go/tools/cmd/staticcheck@latest
go install github.com/stamblerre/gocode@v1.0.0
go install golang.org/x/tools/cmd/goimports@latest

go mod download

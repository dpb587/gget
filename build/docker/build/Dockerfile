FROM golang:1.15.7
RUN apt update && apt install -y xz-utils && rm -rf /var/lib/apt/lists/*
RUN wget -qO- https://github.com/upx/upx/releases/download/v3.96/upx-3.96-amd64_linux.tar.xz | xz -cd | tar -xf- --strip-components=1 -C /usr/local/bin --wildcards '*/upx'

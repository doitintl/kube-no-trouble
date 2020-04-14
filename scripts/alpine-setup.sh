#!/usr/bin/env sh

# Set strict error checking
set -emou pipefail
LC_CTYPE=C

UPX_VERSION="3.96"

apk add --update --no-cache \
	git \
	make \
	tar


wget -qO- "https://github.com/upx/upx/releases/download/v${UPX_VERSION}/upx-${UPX_VERSION}-amd64_linux.tar.xz" \
  | tar --strip 1 -xJv -C "/usr/local/bin/" "upx-${UPX_VERSION}-amd64_linux/upx"

go get github.com/rakyll/statik

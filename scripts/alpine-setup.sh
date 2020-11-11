#!/usr/bin/env sh

# Set strict error checking
set -emou pipefail
LC_CTYPE=C

UPX_VERSION="3.96"
OPA_VERSION="0.22.0"

apk add --update --no-cache \
	curl \
	git \
	make \
	tar


wget -qO- "https://github.com/upx/upx/releases/download/v${UPX_VERSION}/upx-${UPX_VERSION}-amd64_linux.tar.xz" \
  | tar --strip 1 -xJv -C "/usr/local/bin/" "upx-${UPX_VERSION}-amd64_linux/upx"

wget -q -O "/usr/local/bin/opa" "https://github.com/open-policy-agent/opa/releases/download/v${OPA_VERSION}/opa_linux_amd64"
chmod +x "/usr/local/bin/opa"

go get github.com/rakyll/statik
go get github.com/paultyng/changelog-gen

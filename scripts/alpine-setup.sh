#!/usr/bin/env sh

# Set strict error checking
set -emou pipefail
LC_CTYPE=C

OPA_VERSION="0.22.0"

apk add --update --no-cache \
	curl \
	git \
	jq \
	make \
	tar \
	xz

wget -q -O "/usr/local/bin/opa" "https://github.com/open-policy-agent/opa/releases/download/v${OPA_VERSION}/opa_linux_amd64"
chmod +x "/usr/local/bin/opa"

go get github.com/paultyng/changelog-gen

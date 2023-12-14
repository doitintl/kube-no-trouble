#!/usr/bin/env sh

# Set strict error checking
set -eou pipefail
LC_CTYPE=C

OPA_VERSION="0.22.0"

apk add --update --no-cache \
	bash \
	curl \
	git \
	git-cliff \
	jq \
	make \
	tar \
	xz

wget -q -O "/usr/local/bin/opa" "https://github.com/open-policy-agent/opa/releases/download/v${OPA_VERSION}/opa_linux_amd64"
chmod +x "/usr/local/bin/opa"

#!/usr/bin/env sh

# Set strict error checking
set -eou pipefail
LC_CTYPE=C

OPA_VERSION="0.22.0"
CHANGELOG_VERSION="1.1.0"

apk add --update --no-cache \
	bash \
	curl \
	git \
	jq \
	make \
	tar \
	xz

wget -q -O "/usr/local/bin/opa" "https://github.com/open-policy-agent/opa/releases/download/v${OPA_VERSION}/opa_linux_amd64"
chmod +x "/usr/local/bin/opa"

wget -q "https://github.com/paultyng/changelog-gen/releases/download/v${CHANGELOG_VERSION}/changelog-gen_Linux_x86_64.tar.gz" -O - | tar -xz
chmod +x changelog-gen
mv changelog-gen /usr/local/bin/changelog-gen

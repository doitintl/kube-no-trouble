#!/usr/bin/env sh

apk add --update --no-cache \
	make \
	upx

go get github.com/rakyll/statik

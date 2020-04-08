#!/usr/bin/env sh

apk add --update --no-cache \
	make \
	tar \
	upx

go get github.com/rakyll/statik

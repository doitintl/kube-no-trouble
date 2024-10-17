SHELL := bash
.SHELLFLAGS := -eu -o pipefail -c
.ONESHELL:
.SECONDEXPANSION:
.DELETE_ON_ERROR:
.EXPORT_ALL_VARIABLES:
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules

ifeq (, $(shell command -v gtar))
	TAR ?= tar
endif
TAR ?= gtar

ifneq (3.82,$(firstword $(sort $(MAKE_VERSION) 3.82)))
  $(error This Make does not support .ONESHELL, use GNU Make 3.82 and newer)
endif

ifeq (windows,$(GOOS))
  BIN_SUFFIX ?= .exe
endif
BIN_SUFFIX ?=

.DEFAULT_GOAL :=help

GITHUB_REF ?= dev
GIT_REF ?= $(shell echo "$(GITHUB_REF)" | sed "s,refs/[^/]*/,," | tr -cd '[:alnum:]._-')
GITHUB_SHA ?= dev
GIT_SHA ?= $(GITHUB_SHA)

CGO_ENABLED ?= 0
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
BIN_ARCH ?= $(GOOS)-$(GOARCH)
RELEASE_SUFFIX ?= $(GIT_REF)-$(BIN_ARCH).tar.gz

OPA_REPO ?= github.com/open-policy-agent/opa
OPA_VERSION ?= v0.22.0

BIN_DIR ?= bin
CMD_DIR ?= cmd
RELEASE_DIR ?= release-artifacts
CMDS ?= $(shell ls $(CMD_DIR))
BINS ?= $(addsuffix -$(BIN_ARCH)$(BIN_SUFFIX),$(addprefix $(BIN_DIR)/,$(CMDS)))
CHANGELOG ?= changelog.md

RELEASE_ARTIFACTS ?= $(addsuffix -$(RELEASE_SUFFIX),$(addprefix $(RELEASE_DIR)/,$(CMDS)))
SRC ?= $(shell find . -iname '*.go')

GOCMD ?= go
GOBUILD ?= $(GOCMD) build

REQ_BINS = go

_ := $(foreach exec,$(REQ_BINS), \
       $(if $(shell which $(exec)),some string,$(error "No $(exec) binary in $$PATH")))


## Clean, build and pack
all: clean build release-artifacts
.PHONY: all

## Prints list of tasks
help:
	@awk 'BEGIN {FS=":"} /^## .*/,/^[a-zA-Z0-9_-]+:/ { if ($$0 ~ /^## /) { desc=substr($$0, 4) } else { printf "\033[36m%-30s\033[0m %s\n", $$1, desc } }' Makefile | sort
.PHONY: help

## Build binary
build: $(BINS)
.PHONY: build

$(BIN_DIR)/%-$(BIN_ARCH)$(BIN_SUFFIX): $(SRC) go.mod go.sum
	mkdir -p $(BIN_DIR)
	$(GOBUILD) -ldflags="-s -w -X main.version=$(GIT_REF) -X main.gitSha=$(GIT_SHA)" \
	-o "$@" \
	"./$(CMD_DIR)/$(*)"

## Create release artifacts
release-artifacts: $(RELEASE_ARTIFACTS)
.PHONY: release-artifacts

$(RELEASE_DIR)/%-$(RELEASE_SUFFIX): $(BIN_DIR)/%-$(BIN_ARCH)$(BIN_SUFFIX)
	mkdir -p $(RELEASE_DIR)
	$(TAR) -cvz --transform 's,$(BIN_DIR)/$(*)-$(BIN_ARCH)$(BIN_SUFFIX),$(*)$(BIN_SUFFIX),gi' -f "$@" "$<"

## Run Go tests
test: test-fmt
	go test -v -coverprofile fmtcoverage.html ./...
.PHONY: test

## Run go and opt fmt checks
test-fmt:
	test -z "$$(go run $(OPA_REPO)@$(OPA_VERSION) fmt -l pkg/rules/rego/*)"
	test -z "$$(go fmt ./...)"
.PHONY: test-fmt

## Clean build artifacts
clean:
	rm -rf $(BIN_DIR) $(RELEASE_DIR)
.PHONY: clean

## Generate Changelog based on PRs
changelog:
	mkdir -p $(RELEASE_DIR)
	git cliff -u --latest -o $(RELEASE_DIR)/$(CHANGELOG)
.PHONY: changelog

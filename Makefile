GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOARCH=$(shell go env GOARCH)
GOOS=$(shell go env GOOS )

BASE_PATH := $(shell pwd)
BUILD_PATH = $(BASE_PATH)/build

MAIN_PATH=$(BASE_PATH)/main.go
BIN_NAME=mcp-corapanel

.PHONY: build

build:
	mkdir -p $(BUILD_PATH)
	cd $(BASE_PATH) \
	&& CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) $(GOBUILD) -trimpath -ldflags '-s -w' -o $(BUILD_PATH)/$(BIN_NAME) $(MAIN_PATH)

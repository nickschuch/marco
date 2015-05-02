#!/usr/bin/make -f

PROJCET=marco
SHELL=/bin/bash
MKDIR=mkdir
GIT=git
GO=go
RM=rm -rf
CROSS_BASH=source /opt/golang/cross/crosscompile.bash

SOURCE=main.go
TARGETS=darwin-386 darwin-amd64 linux-386 linux-amd64 linux-arm

all: test

build: deps
	@echo "Building..."
	@$(GO) build -o bin/$(PROJCET) $(SOURCE)

deps:
	@echo "Downloading libraries..."
	go-getter Gofile

xbuild: deps dirs
	@for target in $(TARGETS); do \
		echo "Building for $$target..."; \
		$(CROSS_BASH) && \
		$(GO)-$$target build -o bin/$(PROJCET)-$$target $(SOURCE); \
	done;

dirs:
	@$(MKDIR) -p bin

test: build
	@echo "Running tests..."
	@$(GO) test ./...

clean:
	@echo "Cleaning up binaries..."
	$(RM) bin

coverage:
	@echo "Build code coverage..."
	coverage

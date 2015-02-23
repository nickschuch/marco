#!/usr/bin/make -f

SHELL=/bin/bash
MKDIR=mkdir
GIT=git
GO=go
RM=rm -rf
CROSS=https://github.com/davecheney/golang-crosscompile.git
CROSS_BASH=source golang-crosscompile/crosscompile.bash

SOURCE=src/main.go
TARGETS=darwin-386 darwin-amd64 linux-386 linux-amd64 linux-arm

all: test

build: deps
	@echo "Building marco..."
	@$(GO) build -o bin/marco $(SOURCE)

deps:
	@echo "Downloading libraries..."
	@$(GO) get github.com/samalba/dockerclient
	@$(GO) get github.com/Sirupsen/logrus
	@$(GO) get github.com/nickschuch/go-tutum/tutum
	@$(GO) get gopkg.in/alecthomas/kingpin.v1
	@$(GO) get github.com/stretchr/testify/assert

golang-crosscompile:
	$(GIT) clone $(CROSS)
	$(CROSS_BASH) && \
	go-crosscompile-build-all

xbuild: deps golang-crosscompile dirs
	@for target in $(TARGETS); do \
		echo "Building marco for $$target..."; \
		$(CROSS_BASH) && \
		$(GO)-$$target build -o bin/marco-$$target $(SOURCE); \
	done;

dirs:
	@$(MKDIR) -p bin

test: build
	@echo "Run tests..."
	@$(GO) test ./...

clean:
	@echo "Cleanup binaries..."
	$(RM) bin

realclean: clean
	$(RM) golang-crosscompile


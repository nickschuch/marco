#!/usr/bin/make -f

GO=go
GB=gb

all: build

build: clean test
	@echo "Building..."
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GB) build -ldflags '-w -extld ld -extldflags -static'

build-all: build
	@echo "Building others..."
	env GOOS=linux GOARCH=386 $(GB) build
	env GOOS=darwin GOARCH=amd64 $(GB) build
	env GOOS=darwin GOARCH=386 $(GB) build

clean:
	rm -fR pkg bin

test:
	@echo "Running tests..."
	@$(GB) test -test.v=true

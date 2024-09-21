# Makefile for PixivFE

# Variables
BINARY_NAME=pixivfe
TARGETOS ?= $(shell go env GOOS)
TARGETARCH ?= $(shell go env GOARCH)
GIT_COMMIT := $(shell git describe --always 2>/dev/null)

# Include environment variables from .env if it exists
-include .env

.PHONY: all fmt build test run clean install-pre-commit help

all: fmt build test

fmt:
	@echo "Formatting Go code..."
	go fmt ./...

build:
	@echo "Building $(BINARY_NAME)..."
	go mod download
	CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -v -ldflags="-extldflags=-static -X codeberg.org/vnpower/pixivfe/v2/config.REVISION=$(GIT_COMMIT)" -tags netgo -o $(BINARY_NAME)

test:
	@echo "Running tests..."
	go test ./server/template

run: build
	@echo "Running $(BINARY_NAME)..."
	./$(BINARY_NAME)

clean:
	@echo "Cleaning up..."
	rm -f $(BINARY_NAME)

# Additional target to install test script as pre-commit hook
install-pre-commit:
	@echo "Installing pre-commit hook..."
	echo '#!/bin/sh' > .git/hooks/pre-commit
	echo 'go test ./server/template' >> .git/hooks/pre-commit
	chmod +x .git/hooks/pre-commit

# Help target
help:
	@echo "Available targets:"
	@echo "  all                - Run fmt, build, and test"
	@echo "  fmt                - Format Go code"
	@echo "  build              - Build the binary"
	@echo "  test               - Run tests"
	@echo "  run                - Build and run the binary"
	@echo "  clean              - Remove the binary"
	@echo "  install-pre-commit - Install testing pre-commit hook"
	@echo "  help               - Show this help message"

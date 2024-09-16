# Makefile for PixivFE

# Variables
BINARY_NAME=pixivfe
GOFILES=$(shell find . -type f -name '*.go')

# Environment variables (customize as needed)
export PIXIVFE_TOKEN=token_123456
export PIXIVFE_IMAGEPROXY=pximg.cocomi.cf
export PIXIVFE_PORT=8282

.PHONY: all fmt build test run clean

all: fmt build test

fmt:
	@echo "Formatting Go code..."
	@go fmt ./$(shell find . -type d)

build:
	@echo "Building $(BINARY_NAME)..."
	@go mod download
	@go get codeberg.org/vnpower/pixivfe/v2/...
	@CGO_ENABLED=0 go build -v -ldflags="-extldflags=-static" -tags netgo -o $(BINARY_NAME)

test:
	@echo "Running tests..."
	@go test ./server/template

run: build
	@echo "Running $(BINARY_NAME)..."
	./$(BINARY_NAME)

clean:
	@echo "Cleaning up..."
	@rm -f $(BINARY_NAME)

# Additional target to install test script as pre-commit hook
install-pre-commit:
	@echo "Installing pre-commit hook..."
	@cp test.sh .git/hooks/pre-commit
	@chmod +x .git/hooks/pre-commit

# Help target
help:
	@echo "Available targets:"
	@echo "  all              - Run fmt, build, and test"
	@echo "  fmt              - Format Go code"
	@echo "  build            - Build the binary"
	@echo "  test             - Run tests"
	@echo "  run              - Build and run the binary"
	@echo "  clean            - Remove the binary"
	@echo "  install-pre-commit - Install test script as pre-commit hook"
	@echo "  help             - Show this help message"

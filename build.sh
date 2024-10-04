#!/bin/sh

# Variables
BINARY_NAME="pixivfe"
GOOS=${GOOS:-$(go env GOOS)}
GOARCH=${GOARCH:-$(go env GOARCH)}
GIT_COMMIT_DATE=$(git show -s --format=%cd --date=format:"%Y.%m.%d")
GIT_COMMIT_HASH=$(git rev-parse --short HEAD)
REVISION="${GIT_COMMIT_DATE}-${GIT_COMMIT_HASH}"
UNCOMMITTED_CHANGES=$(git status --porcelain)
if [ -n "$UNCOMMITTED_CHANGES" ]; then
    REVISION="${REVISION}+dirty"
fi

fmt() {
    echo "Formatting Go code..."
    go fmt ./...
}

build() {
    echo "Building ${BINARY_NAME}..."
    go mod tidy
    CGO_ENABLED=0 go build -v -ldflags="-extldflags=-static -X codeberg.org/vnpower/pixivfe/v2/config.REVISION=${REVISION}" -o "${BINARY_NAME}"
}

test() {
    echo "Running tests..."
    go test ./test/...
}

run() {
    build
    echo "Running ${BINARY_NAME}..."
    if [ "$1" != "--do-not-load-env-file" ] && [ -f .env ]; then
        echo ".env file found, loading environment variables"
        set -a
        . ./.env
        set +a
    else
        echo "Not loading .env file"
    fi
    ./"${BINARY_NAME}"
}

clean() {
    echo "Cleaning up..."
    rm -f "${BINARY_NAME}"
}

install_pre_commit() {
    echo "Installing pre-commit hook..."
    echo '#!/bin/sh' > .git/hooks/pre-commit
    echo 'go test ./server/template' >> .git/hooks/pre-commit
    chmod +x .git/hooks/pre-commit
}

help() {
    echo "Available commands:"
    echo "  all                - Run fmt, build, and test"
    echo "  fmt                - Format Go code"
    echo "  build              - Build the binary"
    echo "  test               - Run tests"
    echo "  run [--do-not-load-env-file] - Build and run the binary"
    echo "  clean              - Remove the binary"
    echo "  install-pre-commit - Install testing pre-commit hook"
    echo "  help               - Show this help message"
    echo ""
    echo "Options:"
    echo "  --do-not-load-env-file - Do not load the .env file when running"
}

all() {
    fmt
    build
    test
}

# Function to handle command execution
execute_command() {
    case "$1" in
        fmt) fmt ;;
        build) build ;;
        test) test ;;
        run)
            if [ "$2" = "--do-not-load-env-file" ]; then
                run "--do-not-load-env-file"
            else
                run
            fi
            ;;
        clean) clean ;;
        install-pre-commit) install_pre_commit ;;
        help) help ;;
        all) all ;;
        *)
            echo "Unknown command: $1"
            echo "Use 'help' to see available commands"
            exit 1
            ;;
    esac
}

# Main execution
if [ $# -eq 0 ]; then
    build
else
    execute_command "$@"
fi

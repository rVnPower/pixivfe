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
    go test ./...
}

scan() {
    semgrep scan -q -f semgrep.yml
}

i18n_code() {
    semgrep scan -q -f i18n/semgrep-i18n.yml --json | jq '.results | map({msg:.extra.metavars["$MSG"].abstract_content, file:.path, line:.start.line, offset:.start.offset})' > i18n/code_strings.json
}

i18n_template() {
    go run ./i18n/crawler > i18n/template_strings.json
    echo "Malformed strings are listed below:"
    jq '.[] | select(.msg | contains("\n"))' < i18n/template_strings.json
}

i18n() {
    i18n_code
    i18n_template
    mkdir -p i18n/locale/en
    go run ./i18n/converter < i18n/code_strings.json > i18n/locale/en/code.json
    go run ./i18n/converter < i18n/template_strings.json > i18n/locale/en/template.json
    chmod -w i18n/locale/en/code.json i18n/locale/en/template.json
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
    echo './build.sh test' >> .git/hooks/pre-commit
    chmod +x .git/hooks/pre-commit
}

help() {
    echo "Available commands:"
    echo "  all                - Run fmt, build, and test"
    echo "  fmt                - Format Go code"
    echo "  build              - Build the binary"
    echo "  scan               - Scan Go code"
    echo "  test               - Run tests"
    echo "  i18n               - Extract i18n strings"
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
    i18n
    build
    test
}

# Function to handle command execution
execute_command() {
    case "$1" in
        fmt) fmt ;;
        build) build ;;
        test) test ;;
        scan) scan ;;
        i18n) i18n ;;
        i18n_code) i18n_code ;;
        i18n_template) i18n_template ;;
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

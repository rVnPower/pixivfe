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

# Check for .env file and load it unless explicitly told not to
if [ "$1" != "--no-env-file" ] && [ -f .env ]; then
    echo ".env file found, loading environment variables"
    set -a
    . ./.env
    set +a
else
    echo "Not loading .env file"
fi

fmt() {
    echo "Formatting Go code..."
    go fmt ./...
}

build() {
    echo "Building ${BINARY_NAME}..."
    go mod tidy
    i18n
    CGO_ENABLED=0 go build -v -ldflags="-extldflags=-static -X codeberg.org/vnpower/pixivfe/v2/config.REVISION=${REVISION}" -o "${BINARY_NAME}"
}

build_docker() {
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
    echo "Extracting i18n strings from code..."
    mkdir -p i18n/locale/en
    rm -f i18n/locale/en/code.json
    semgrep scan -q -f i18n/semgrep-i18n.yml --json | jq '.results | map({msg:.extra.metavars["$MSG"].abstract_content, file:.path, line:.start.line, offset:.start.offset})' > i18n/code_strings.json
    go run ./i18n/converter < i18n/code_strings.json > i18n/locale/en/code.json
    chmod -w i18n/locale/en/code.json
}

i18n_template() {
    echo "Extracting i18n strings from templates..."
    mkdir -p i18n/locale/en
    rm -f i18n/locale/en/template.json
    go run ./i18n/crawler > i18n/template_strings.json
    go run ./i18n/converter < i18n/template_strings.json > i18n/locale/en/template.json
    chmod -w i18n/locale/en/template.json
    malformed_strings=$(jq 'to_entries | .[] | select(.value | contains("\n"))' < i18n/locale/en/template.json)
    if [ -z "$malformed_strings" ]; then
        echo "No malformed strings found."
    else
        echo "Malformed strings are listed below:"
        echo "$malformed_strings"
    fi
}

i18n() {
    echo "Starting i18n extraction process..."
    i18n_code
    i18n_template
    echo "i18n extraction completed."
}

i18n_upload() {
    echo "Uploading i18n strings to Crowdin..."
    crowdin upload
}

i18n_download() {
    echo "Downloading i18n strings from Crowdin..."
    crowdin download
}

run() {
    build
    echo "Running ${BINARY_NAME}..."
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
    echo "Pre-commit hook installed successfully."
}

help() {
    echo "Available commands:"
    echo "  all                - Run fmt, build, and test"
    echo "  fmt                - Format Go code"
    echo "  build              - Build the binary"
    echo "  build_docker       - Build the binary (for Docker, skips i18n refresh)"
    echo "  scan               - Scan Go code"
    echo "  test               - Run tests"
    echo "  i18n               - Extract i18n strings"
    echo "  i18n-up            - Upload strings to Crowdin"
    echo "  i18n-down          - Download strings from Crowdin"
    echo "  run                - Build and run the binary"
    echo "  clean              - Remove the binary"
    echo "  install-pre-commit - Install testing pre-commit hook"
    echo "  help               - Show this help message"
    echo ""
    echo "Options:"
    echo "  --no-env-file - Do not load the .env file (must be the first argument if used)"
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
        build_docker) build_docker ;;
        test) test ;;
        scan) scan ;;
        i18n) i18n ;;
        i18n_code) i18n_code ;;
        i18n_template) i18n_template ;;
        i18n-up) i18n_upload ;;
        i18n_upload) i18n_upload ;;
        i18n-down) i18n_download ;;
        i18n_download) i18n_download ;;
        run) run ;;
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
elif [ "$1" = "--no-env-file" ]; then
    shift
    if [ $# -eq 0 ]; then
        build
    else
        execute_command "$@"
    fi
else
    execute_command "$@"
fi

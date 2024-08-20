#!/bin/sh

go test ./routes

# install as pre-commit hook:
# cp test.sh ~/.git/hooks/pre-commit
#!/bin/bash

PROJ_DIR=$1
cd "$PROJ_DIR" || exit

source ./scripts/prompt.sh

prompt "Running lint"

golangci-lint run ./...
govulncheck ./...

end

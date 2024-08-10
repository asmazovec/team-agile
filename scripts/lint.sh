#!/bin/bash

SCRIPT_DIR="$(dirname "$(readlink -f "$0")")"
source "$SCRIPT_DIR/prompt.sh"

prompt "Running lint"

golangci-lint run ./...
govulncheck ./...

end

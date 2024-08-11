#!/bin/bash

PROJ_DIR=$1
cd "$PROJ_DIR" || exit

source ./scripts/prompt.sh

STATUS=0

prompt "Running golangci-lint"
if ! golangci-lint run ./...; then STATUS=1; fi

prompt "Running govulncheck"
if ! govulncheck ./...; then STATUS=1; fi

end
exit "$STATUS"

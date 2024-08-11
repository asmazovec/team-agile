#!/bin/bash

PROJ_DIR=$1
cd "$PROJ_DIR" || exit 1

source ./scripts/prompt.sh

prompt "Restoring packages"

go mod tidy
go mod verify

end

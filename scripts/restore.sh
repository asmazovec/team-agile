#!/bin/bash

SCRIPT_DIR="$(dirname "$(readlink -f "$0")")"
source "$SCRIPT_DIR/prompt.sh"

prompt "Restoring packages"

go mod tidy
go mod verify

end

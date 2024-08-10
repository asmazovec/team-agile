#!/bin/bash

SCRIPT_DIR="$(dirname "$(readlink -f "$0")")"
source "$SCRIPT_DIR/prompt.sh"

# See CONTRIBUTING.md for details
if ! command -v direnv &> /dev/null
then
  prompt "! Skipping direnv... It is recommended to direnv !"
else
  direnv allow
fi

prompt "Restoring tools packages"

cd internal/tools
go mod tidy
go mod verify

end

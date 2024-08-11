#!/bin/bash

PROJ_DIR=$1
cd "$PROJ_DIR" || exit 1

source ./scripts/prompt.sh

# See CONTRIBUTING.md for details
if ! command -v direnv &> /dev/null
then
  prompt "! Skipping direnv... It is recommended to direnv !"
else
  direnv allow
fi

prompt "Restoring tools packages"

cd ./internal/tools || exit 1
go mod tidy
go mod verify

end

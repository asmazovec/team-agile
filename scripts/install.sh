#!/bin/bash

PROJ_DIR=$1
cd "$PROJ_DIR" || exit 1

source ./scripts/prompt.sh

GO_VERSION=$2

# See CONTRIBUTING.md for details
prompt "Getting Go SDK"

go install "golang.org/dl/go$GO_VERSION@latest"
go"$GO_VERSION" download
mkdir -p bin
go get ./...
ln -sf "$(go env GOPATH)/bin/go$GO_VERSION" bin/go

end

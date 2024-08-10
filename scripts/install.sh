#!/bin/sh

GO_VERSION=$1

SCRIPT_DIR="$(dirname "$(readlink -f "$0")")"
source "$SCRIPT_DIR/prompt.sh"

# See CONTRIBUTING.md for details
echo `prompt "Getting Go SDK"`

go install golang.org/dl/go${GO_VERSION}@latest
go${GO_VERSION} download
mkdir -p bin
go get ./...
ln -sf `go env GOPATH`/bin/go${GO_VERSION} bin/go

end

#!/bin/sh

GO_VERSION=$1

MYDIR="$(dirname "$(readlink -f "$0")")"
source $MYDIR"/prompt.sh"

# Getting Go SDK
# See CONTRIBUTING.md for details
echo `prompt "Getting Go SDK"`
go install golang.org/dl/go${GO_VERSION}@latest
go${GO_VERSION} download
mkdir -p bin
go get ./...
ln -sf `go env GOPATH`/bin/go${GO_VERSION} bin/go

# Enable direnv utility
# See CONTRIBUTING.md for details
if ! command -v direnv &> /dev/null
then
  echo `prompt "Skipping direnv... It is recommended to direnv!"`
else
  echo `prompt "Enabling direnv"`
  direnv allow
  echo "done"
fi

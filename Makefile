.DEFAULT_GOAL = install
.PHONY: FORCE
PROJ_DIR="$(dirname "$(readlink -f "$0")")"

export GOPROXY = https://proxy.golang.org
export GO_VERSION = 1.22.5


######### Actions #########

.PHONY: lint
lint: install/tools
	./scripts/lint.sh

.PHONY: install
install: go.sum
	./scripts/install.sh ${GO_VERSION}

.PHONY: install/dev
install/dev: install install/tools

.PHONY: install/tools
install/tools:
	./scripts/tools.sh

.PHONY: test
test:
	go test -v -shuffle=on -race -coverprofile=docs/coverage.txt -covermode=atomic `go list ./... | grep -v /cmd/`

.PHONY: docs/coverage.html
docs/coverage: docs/coverage.html


######### Static #########

docs/coverage.html:
	go tool cover -html docs/coverage.txt -o docs/coverage.html

docs/:
	mkdir -p docs

go.sum: go.mod

go.mod: FORCE
	./scripts/restore.sh

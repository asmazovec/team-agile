.DEFAULT_GOAL = install
.PHONY: FORCE
PROJ_DIR="$(dirname "$(readlink -f "$0")")"

export GOPROXY = https://proxy.golang.org
export GO_VERSION = 1.22.5

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

go.mod: FORCE
	./scripts/restore.sh

go.sum: go.mod

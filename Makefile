.DEFAULT_GOAL = install
.PHONY: FORCE

export GOPROXY = https://proxy.golang.org
export GO_VERSION = 1.22.5

MAKEFILE := $(abspath $(lastword $(MAKEFILE_LIST)))
PROJ_DIR := $(patsubst %/,%,$(dir $(MAKEFILE)))

TEST_SKIP := /cmd/
TEST_DIRS := $$(go list ${PROJ_DIR}/... | grep -v ${TEST_SKIP})
TEST_OPTS := -shuffle=on -race -covermode=atomic -coverprofile=${TEST_COVER_OUT}

TEST_COVER_OUT := ${PROJ_DIR}/docs/coverage.txt
TEST_COVER_CMD := go tool cover

######### Actions #########

.PHONY: install
install: go.sum
	${PROJ_DIR}/scripts/install.sh ${PROJ_DIR} ${GO_VERSION}

.PHONY: install/tools
install/tools: install
	${PROJ_DIR}/scripts/tools.sh ${PROJ_DIR}

.PHONY: install/dev
install/dev: install install/tools
	git config core.hooksPath ${PROJ_DIR}/scripts/hooks
	go install -C ${PROJ_DIR}/internal/tools \
		github.com/golangci/golangci-lint/cmd/golangci-lint \
		golang.org/x/tools/cmd/goimports \
		golang.org/x/vuln/cmd/govulncheck

.PHONY: lint
lint:
	${PROJ_DIR}/scripts/lint.sh ${PROJ_DIR}

.PHONY: test
test:
	go test ${TEST_OPTS} ${TEST_DIRS}

.PHONY:
test/coverage:
	go tool cover -func=${TEST_COVER_OUT}

.PHONY: docs/coverage.html
docs/coverage: docs/coverage.html


######### Static #########

docs/coverage.html:
	${TEST_COVER_CMD} -html ${PROJ_DIR}/docs/coverage.txt -o ${PROJ_DIR}/docs/coverage.html

docs/:
	mkdir -p docs

go.sum: go.mod

go.mod: FORCE
	${PROJ_DIR}/scripts/restore.sh ${PROJ_DIR}

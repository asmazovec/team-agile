.DEFAULT_GOAL = install
.PHONY: FORCE

export GOPROXY = https://proxy.golang.org
export GO_VERSION = 1.22.5

MAKEFILE    := $(abspath $(lastword $(MAKEFILE_LIST)))
PROJ_DIR    := $(patsubst %/,%,$(dir $(MAKEFILE)))
DOCS_DIR    := ${PROJ_DIR}/docs
SCRIPTS_DIR := ${PROJ_DIR}/scripts

TEST_COVER_CONFIG  := ${PROJ_DIR}/.testcoverage.yml
TEST_COVER_PROFILE := ${DOCS_DIR}/coverage.txt
TEST_COVER_DOCS    := ${DOCS_DIR}/coverage.html
TEST_COVER_OPTS    := -p ${TEST_COVER_PROFILE} -c ${TEST_COVER_CONFIG}

TEST_SKIP := /cmd/
TEST_DIRS := $$(go list ${PROJ_DIR}/... | grep -v ${TEST_SKIP})
TEST_OPTS := -shuffle=on -race -covermode=atomic -coverprofile=${TEST_COVER_PROFILE}

######### ACTIONS #########

.PHONY: clean
clean:
	rm -f ${TEST_COVER_PROFILE}
	rm -f ${TEST_COVER_DOCS}

.PHONY: install
install: go.sum
	${SCRIPTS_DIR}/install.sh ${PROJ_DIR} ${GO_VERSION}

.PHONY: install/tools
install/tools: install
	${SCRIPTS_DIR}/tools.sh ${PROJ_DIR}
	go install -C ${PROJ_DIR}/internal/tools \
		github.com/golangci/golangci-lint/cmd/golangci-lint \
		golang.org/x/tools/cmd/goimports \
		golang.org/x/vuln/cmd/govulncheck \
		github.com/vladopajic/go-test-coverage/v2

.PHONY: install/dev
install/dev: install install/tools
	git config core.hooksPath ${SCRIPTS_DIR}/hooks


######### CHECKS #########

.PHONY: check/lint
check/lint:
	golangci-lint run ./...

.PHONY: check/security
check/security:
	govulncheck ./...

.PHONY: check/test
check/test:
	rm -f ${TEST_COVER_PROFILE}
	go test ${TEST_OPTS} ${TEST_DIRS}

.PHONY: check/coverage
check/coverage: check/test
	go-test-coverage ${TEST_COVER_OPTS}


######### DOCS #########

docs/coverage.html: check/test docs/
	go tool cover -html ${TEST_COVER_PROFILE} -o ${DOCS_DIR}/coverage.html


##########################

docs/:
	mkdir -p ${DOCS_DIR}

go.sum: go.mod

go.mod: FORCE
	${SCRIPTS_DIR}/restore.sh ${PROJ_DIR}

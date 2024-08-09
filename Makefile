.DEFAULT_GOAL = test
.PHONY: FORCE

export GOPROXY = https://proxy.golang.org
export GO_VERSION = 1.22.5

.PHONY: lint
lint: install/tools tool/golangci-lint tool/govulncheck tool/vet

.PHONY: install/tools
install/tools:
	go install -C internal/tools \
		github.com/golangci/golangci-lint/cmd/golangci-lint \
		golang.org/x/vuln/cmd/govulncheck

.PHONY: install
install:
	go install golang.org/dl/go${GO_VERSION}@latest
	go${GO_VERSION} download
	mkdir -p bin
	go get ./...
	ln -sf `go env GOPATH`/bin/go${GO_VERSION} bin/go

.PHONY: tool/golangci-lint
tool/golangci-lint:
	golangci-lint run

.PHONY: tool/govulncheck
tool/govulncheck:
	govulncheck ./...

.PHONY: tool/vet
tool/vet:
	go vet ./...

go.mod: FORCE
	go mod tidy
	go mod verify

go.sum: go.mod


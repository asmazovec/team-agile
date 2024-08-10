# Contributing to the project

To contribute to this project, please follow these recommendations.

## Pre requirements

To contribute to this project, your system should be prepared with some utilities:

**required**
- `go`  https://go.dev/ - this project uses the Go language and its powerful toolset.
- `make` https://www.gnu.org/software/make/ - the make build system allows you to call functional script shortcuts.

**recommended**
- `direnv` https://direnv.net/ - lets establish project-specific environment variable separately from the operating system. 
See also "[Approach to Manage Project Tools Versions](internal/tools/README.md)".


## Setup project for developing

Project setup works around the make build system.

```bash
$ make install        # install go with current versions
$ make install/tools  # install go tools
```

## Lint check

The project's source code should follow the restrictions and code style conventions established for this project.
The following commands allow you to perform automated verification with static check analysis and alert you to some code vulnerabilities.

```bash
# main verification shortcut
$ make lint

# additional options
$ make tool/golangci-lint
$ make tool/vet
$ make tool/govulncheck
```
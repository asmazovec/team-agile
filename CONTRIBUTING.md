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
# common approach
$ make install/dev    # prepare project for development

# additional options (included in common approach)
$ make install        # install project with go and dependencies
$ make install/tools  # install project tools
```

## Checks

The project's source code should follow the restrictions and code style conventions established for this project.
The following commands allow you to perform automated verification with static check analysis and alert you to some code vulnerabilities.

```bash
$ make check/lint     # check for code style and whatever static checks
$ make check/security # check for security vulnerabilities
```

## Tests

Run tests with correct flags and check coverage Run tests with correct flags and check coverage these commands are used.

```bash
$ make check/test         # run tests, also generates coverage profile along the project.
$ make check/coverage     # check for satisfying test coverage thresholds
$ make docs/coverage.html # generate coverage info docs
```

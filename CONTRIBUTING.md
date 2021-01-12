# Implementation details
To streamline dev efforts, the opctl project is maintained as a monorepo containing semi-independent sub-projects.

Typically at a minimum, each sub-project includes it's own `README.md`, `CONTRIBUTING.md`, and`.opspec`. 

The project is configured as a single go module, which allows dependencies to be shared during development, and should allow most IDEs to understand the project configuration without manual configuration.

# Project structure

## [api](api)
OpenAPI spec for the opctl ReST API.

## [cli](cli)
CLI for opctl.

## [opspec](opspec)
JSON schema for the opspec language.

used by:
- [sdks/go](sdks/go)

## [sdks/go](sdks/go)
SDK for integrating with opctl from golang.

used by:
- [cli](cli)

## [sdks/js](sdks/js)
SDK (written in typescript) for integrating with opctl from javascript/typescript.

used by:
- [webapp](webapp)
- [sdks/react](sdks/react)

## [sdks/react](sdks/react)
SDK (written in typescript) for integrating with opctl from react.

used by:
- [webapp](webapp)

## [test-suite](test-suite)
End to end tests, maintained as ops, for testing opctl and SDKs.

used by:
- [cli](cli)
- [sdks/go](sdks/go)

## [webapp](webapp)
Webapp for opctl.

used by:
- [cli](cli)

## [website](website)
Website hosted at [https://opctl.io](https://opctl.io).

# Pull requests
Pull requests are subject to:

- approval by one or more [maintainers](https://github.com/orgs/opctl/teams/maintainers/members)
- the [build](.opspec/build) op continuing to run with a successful outcome

# Local development with native go tooling

First, ensure dependencies are installed

- [go](https://golang.org/doc/install) (use the version from the [go.mod file](./go.mod#L3))
- [gpgme](https://www.gnupg.org/related_software/gpgme/) (on macOS, `brew install gpgme`)
- [dlv](https://github.com/go-delve/delve) (for debugging) - Because this project uses go modules, install this globally with `go get github.com/go-delve/delve` to run the tool outside of the project directory.

Build a binary with

`go build -o opctl-beta ./cli`

To debug

`go run github.com/go-delve/delve/cmd/dlv --check-go-version=false --listen=127.0.0.1:40000 --headless=true --api-version=2 exec /PATH/TO/OPCTL/opctl-beta run dev`

This will suspend execution until a client connects on port 40000. If using VSCode, you can use the saved run configuration to connect.

On macOS, you can debug exit signal handling by sending an interrupt to the underlying PID. When the debugger connects, it will give you the with a message like: `Got a connection, launched process ../opctl/opctl-beta (pid = 79924).` (path to process and PID will be different). Kill with `kill -SIGINT 72958`.

# Implementation details
The CLI is built using [mow](https://github.com/jawher/mow.cli) and tries to rely on the [opctl go SDK](../sdks/go/README.md) whenever behavior is not specific to the CLI.

The [opctl webapp](../webapp/README.md) is statically embedded and hosted at runtime.


# How do I...

## List operations
1. `opctl ls` from this directory will print a full operation manual.

## Build
1. `opctl run compile` from this directory to compile and test the CLI. On success, `opctl-linux-amd64`, `opctl-darwin-amd64` will exist in this directory.

## (Re)Vendor
1. `go mod vendor` from this directory will (re)vendor all go dependencies.

## Pull in local [go SDK](../sdks/go/README.md) changes
1. same as [(Re)Vendor](#revendor).

## Test
1. `opctl run test` from this directory to compile and test the CLI.

## (Re)Generate fake implementations
1. `opctl run generate` from this directory to (re)generate fake implementations of interfaces as they're added/changed.


# Contribution guidelines
- DO follow [go code review comments](https://github.com/golang/go/wiki/CodeReviewComments).
- DO write tests in `arrange`, `act`, `assert` format w/ the given object under test referred to as `objectUnderTest`.
- DO keep tests alongside source code; i.e. place `code_test.go` alongside `code.go`.
- DO depend on interfaces, not implementations, and use [fakes](https://github.com/maxbrunsfeld/counterfeiter) to test dependency interactions.

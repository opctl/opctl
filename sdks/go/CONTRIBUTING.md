# Implementation details
The go SDK is maintained as a go module.

# How do I...

## List operations
1. `opctl ls` from this directory will print a full operation manual.

## (Re)Vendor
1. `go mod vendor` from this directory will (re)vendor all go dependencies.

## Test
1. `opctl run test` from this directory to compile and test the SDK.

## (Re)Generate fake implementations
1. `opctl run generate` from this directory to (re)generate fake implementations of interfaces as they're added/changed.


# Contribution guidelines
- DO follow [go code review comments](https://github.com/golang/go/wiki/CodeReviewComments).
- DO write tests in `arrange`, `act`, `assert` format w/ the given object under test referred to as `objectUnderTest`.
- DO keep tests alongside source code; i.e. place `code_test.go` alongside `code.go`.
- DO depend on interfaces, not implementations, and use [fakes](https://github.com/maxbrunsfeld/counterfeiter) to test dependency interactions.

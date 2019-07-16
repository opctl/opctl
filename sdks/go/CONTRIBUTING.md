# Dev ops

Ops are maintained in
[![opspec 0.1.6](https://img.shields.io/badge/opspec-0.1.6-brightgreen.svg?colorA=6b6b6b&colorB=fc16be)](https://opspec.io)
definition format.

They can be consumed via tools like [opctl](https://opctl.io).

# Acceptance criteria

Contributions are subject to:

- accepted review by one or more
  [maintainers](https://github.com/orgs/opctl/teams/maintainers/members)
- the [build](.opspec/build) op continuing to run with a successful
  outcome
- adherence to
  [go code review comments](https://github.com/golang/go/wiki/CodeReviewComments)
  
# Dependency management

 [dep](https://golang.github.io/dep/) is used to manage go dependencies
All go dependencies are vendored. 

# Documentation

## GoDoc

Documentation for SDK packages are maintained in golang's native go doc format; which is web browsable via the [godoc webpage](http://godoc.org/github.com/opctl/sdk-golang)


# Testing

`opctl run test` runs all tests inclusive of code coverage.

## Fakes

To streamline unit test related maintenance, [counterfeiter](https://github.com/maxbrunsfeld/counterfeiter) is used to auto-generate fake implementations of interfaces.

The fakes are then used to assert on & stub the object under tests interactions w/ its dependencies. 

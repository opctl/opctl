# Dev ops

Ops are maintained in
[![opspec 0.1.5](https://img.shields.io/badge/opspec-0.1.5-brightgreen.svg?colorA=6b6b6b&colorB=fc16be)](https://opspec.io/0.1.5/packages.html)
format.

They can be consumed via tools like [opctl](https://opctl.io).

# Acceptance criteria

Contributions are subject to:

- accepted review by one or more
  [maintainers](https://github.com/orgs/opspec-io/teams/maintainers/members)
- the [build](.opspec/build) op continuing to run with a successful
  outcome
- adherence to
  [go code review comments](https://github.com/golang/go/wiki/CodeReviewComments)


# Testing

`opctl run test` runs all tests inclusive of code coverage.

## Fakes

To streamline unit test related maintenance, [counterfeiter](https://github.com/maxbrunsfeld/counterfeiter) is used to auto-generate fake implementations of interfaces.

The fakes are then used to assert on & stub the object under tests interactions w/ its dependencies. 
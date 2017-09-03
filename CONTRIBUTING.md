# Dev ops

Ops are maintained in
[![opspec 0.1.4](https://img.shields.io/badge/opspec-0.1.4-brightgreen.svg)](https://opspec.io/0.1.4/packages.html#format)
package format.

They can be consumed via tools like [opctl](https://opspec.io/opctl).

# Repo layout

## [/vendor](vendor)

root path for external deps of this repo. This includes vendored go pkgs
& git sub modules.

# Acceptance criteria

Contributions are subject to:

- accepted review by a  by 66% of the projects maintainers (see
  [MAINTAINERS.md](MAINTAINERS.md))
- the [build](.opspec/build) op continuing to run with a successful
  outcome
- adherence to
  [go code review comments](https://github.com/golang/go/wiki/CodeReviewComments)

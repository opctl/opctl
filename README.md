[![Build Status](https://travis-ci.org/opspec-io/opctl.svg?branch=master)](https://travis-ci.org/opspec-io/opctl)
[![codecov](https://codecov.io/gh/opspec-io/opctl/branch/master/graph/badge.svg)](https://codecov.io/gh/opspec-io/opctl)

[opspec](https://opspec.io) compliant runtime

> *Be advised: this project is currently at Major version zero. Per the
> semantic versioning spec: "Major version zero (0.y.z) is for initial
> development. Anything may change at any time. The public API should
> not be considered stable."*

# Usage

for usage guidance simply execute without any arguments:

```SHELL
opctl

Usage: opctl [OPTIONS] COMMAND [arg...]

https://opspec.io compliant runtime

Options:
  -v, --version            Show the version and exit
  --nc, --no-color=false   Disable output coloring

Commands:
  collection    Collection related actions
  daemon        Run the opctl daemon
  events        Stream events
  kill          Kill an op
  ls            List ops in a collection
  op            Op related actions
  run           Run an op
  self-update   Update opctl

Run 'opctl COMMAND --help' for more information on a command.
```

# Releases

for every release:
- source code will be [tagged](https://github.com/opspec-io/opctl/tags).
- platform specific binaries/installers will be made available on [opctl.io](https://opctl.io)

# Versioning

This project adheres to the [Semantic Versioning](http://semver.org/)
specification

# Contributing

see [CONTRIBUTING.md](CONTRIBUTING.md)

# Changelog

see [CHANGELOG.md](CHANGELOG.md)

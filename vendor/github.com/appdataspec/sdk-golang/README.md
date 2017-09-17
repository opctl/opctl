[![Build Status](https://travis-ci.org/appdataspec/sdk-golang.svg?branch=master)](https://travis-ci.org/appdataspec/sdk-golang)[![Coverage](https://codecov.io/gh/appdataspec/sdk-golang/branch/master/graph/badge.svg)](https://codecov.io/gh/appdataspec/sdk-golang)

Golang SDK for the [app data spec](https://github.com/appdataspec/spec)

> *Be advised: this project is currently at Major version zero. Per the
> semantic versioning spec: "Major version zero (0.y.z) is for initial
> development. Anything may change at any time. The public API should
> not be considered stable."*

# Usage

```go
package myDummyPackage

import (
"github.com/appdataspec/sdk-golang/appdatapath"
"fmt"
)

func main() {
// use path package for working w/ spec compliant app data paths
appDataPath := appdatapath.New()

fmt.Printf("Global path is: %v\n", appDataPath.Global())
fmt.Printf("Per user path is: %v\n", appDataPath.PerUser())
}
```

# Releases

All releases will be
[tagged](https://github.com/appdataspec/sdk-golang/tags) and made
available on the
[releases](https://github.com/appdataspec/sdk-golang/releases) page with
links to docs.

# Versioning

This project adheres to the [Semantic Versioning](http://semver.org/)
specification

# Contributing

see [CONTRIBUTING.md](CONTRIBUTING.md)

# Changelog

see [CHANGELOG.md](CHANGELOG.md)

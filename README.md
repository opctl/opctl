[![Build Status](https://travis-ci.org/opspec-io/engine.svg?branch=master)](https://travis-ci.org/opspec-io/engine)
[![codecov](https://codecov.io/gh/opspec-io/engine/branch/master/graph/badge.svg)](https://codecov.io/gh/opspec-io/engine)

[opspec](http://opspec.io) compliant engine.

*Be advised: this project is currently at Major version zero. Per the
semantic versioning spec: "Major version zero (0.y.z) is for initial
development. Anything may change at any time. The public API should not
be considered stable."*

# Official SDK's

[engine-sdk-golang](https://github.com/opspec-io/engine-sdk-golang)

# Supported Use Cases

- get event stream
- get liveness
- kill op run
- run op

# Example Usage

### 1) Start dockerized engine

```SHELL
docker run \
-d \ # run as daemon
-v /var/run/docker.sock:/var/run/docker.sock \ # allow access to docker socket
-v /Users:/Users \ # enable ops in host `/Users` dir
-p 42224:42224 \ # expose engine TCP API
--network host \ # use host network (performance enhancement)
--restart always \ # restart on failure
--name opctl_engine \
opspec/engine # unstable version of engine
```

### 2) Explore the engine API via Swagger UI

**On host with docker-machine**:  
open your browser and navigate to the url returned from `echo
$(docker-machine ip):42224`

**On host with native docker**:  
open your browser and navigate to 127.0.0.1:42224

# Releases

All releases will be [tagged](https://github.com/opspec-io/engine/tags) and
made available on the [releases](https://github.com/opspec-io/engine/releases)
page with links to the corresponding versions of the
[INSTALLATION.md](INSTALLATION.md) and [CHANGELOG.md](CHANGELOG.md)
docs.

# Versioning

This project adheres to the [Semantic Versioning](http://semver.org/)
specification

# Contributing

see [CONTRIBUTING.md](CONTRIBUTING.md)

# Changelog

see [CHANGELOG.md](CHANGELOG.md)

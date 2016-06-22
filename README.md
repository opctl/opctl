[![Build Status](https://travis-ci.org/opctl/engine.svg?branch=master)](https://travis-ci.org/opctl/engine)
[![codecov](https://codecov.io/gh/opctl/engine/branch/master/graph/badge.svg)](https://codecov.io/gh/opctl/engine)

engine for running http://opspec.io compliant ops.

*Be advised: this project is currently at Major version zero. Per the
semantic versioning spec: "Major version zero (0.y.z) is for initial
development. Anything may change at any time. The public API should not
be considered stable."*

# Official SDK's

[engine-sdk-golang](https://github.com/opctl/engine-sdk-golang)

# Supported Use Cases

- get event stream
- get liveness
- kill op run
- run op

# Example Usage

### 1) Start dockerized engine

```SHELL
docker run -it --rm -v /var/run/docker.sock:/var/run/docker.sock -v /Users:/Users -p 42224:42224 opctl/engine:0.1.1
```

explanation:

- `-it` interactive/tty
- `--rm` remove on exit
- `-v /var/run/docker.sock:/var/run/docker.sock` bind mount host docker
  socket
- `-v /Users:/Users` bind mount host `/Users` dir
- `-p 42224:42224` expose container port `42224` via docker-machine port
  `42224`
- `opctl/engine:0.1.1` use latest
  [opctl/engine](https://hub.docker.com/r/opctl/engine/) image

### 2) Explore the engine API via Swagger UI

**On host with docker-machine**:  
open your browser and navigate to the url returned from `echo
$(docker-machine ip):42224`

**On host with native docker**:  
open your browser and navigate to 127.0.0.1:42224

# Releases

All releases will be [tagged](https://github.com/opctl/engine/tags) and
made available on the
[releases](https://github.com/opctl/engine/releases) page with release
notes.

# Versioning

This project adheres to the [Semantic Versioning](http://semver.org/)
specification

# Contributing

see [CONTRIBUTING.md](CONTRIBUTING.md)

# Changelog

see [CHANGELOG.md](CHANGELOG.md)

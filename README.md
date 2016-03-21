[![Build Status](https://travis-ci.org/dev-op-spec/engine.svg?branch=master)](https://travis-ci.org/dev-op-spec/engine)
[![codecov.io](https://codecov.io/github/dev-op-spec/engine/coverage.svg?branch=master)](https://codecov.io/github/dev-op-spec/engine?branch=master)

An engine for the dev op spec

# Official SDK's

[sdk-for-golang](https://github.com/dev-op-spec/sdk-for-golang)

# Supported Use Cases

### Dev Ops
- add dev op
- list dev ops
- run dev op
- set description of dev op

### Pipelines
- add stage to pipeline
- add pipeline
- list pipelines
- run pipeline
- set description of pipeline

# Runtime Dependencies

The environment in which the sdk executes, must have the following available on its $Path for 
full functionality:

- [docker](https://github.com/docker/docker) >= 1.10
- [docker-compose](https://github.com/docker/compose) >= 1.6
- [git](https://git-scm.com/) >= 1.8.0

Note: if using Windows or OSX, you need to update your docker-machine to use NFS instead of vboxfs 
(or suffer painfully slow performance). One recommended way to achieve this is via 
[docker-machine-nfs](https://github.com/adlogix/docker-machine-nfs). 
Your mileage may vary.

# Installation

see [releases](https://github.com/dev-op-spec/engine/releases) for available versions and release notes

# Contributing

refer to [CONTRIBUTING.md](CONTRIBUTING.md)

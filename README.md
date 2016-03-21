[![Build Status](https://travis-ci.org/dev-op-spec/engine.svg?branch=master)](https://travis-ci.org/dev-op-spec/engine)
[![codecov.io](https://codecov.io/github/dev-op-spec/engine/coverage.svg?branch=master)](https://codecov.io/github/dev-op-spec/engine?branch=master)

A lightweight runtime for interacting with dev op specs.

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

# Prerequisites

- [docker](https://github.com/docker/docker) >= 1.10

Note: if using Windows or OSX, you need to update your docker-machine to use NFS instead of vboxfs 
(or suffer painfully slow performance). One recommended way to achieve this is via 
[docker-machine-nfs](https://github.com/adlogix/docker-machine-nfs). 
Your mileage may vary.

# Example Usage

### 1) Start engine for current project
```SHELL
docker run -it --rm -v /var/run/docker.sock:/var/run/docker.sock -v $(pwd):$(pwd) -w $(pwd) -p 8080:8080 devopspec/engine
```
explanation:

- `-it` interactive/tty
- `--rm` remove on exit
- `-v /var/run/docker.sock:/var/run/docker.sock` bind mount host docker socket
- `-v $(pwd):$(pwd)` bind mount host workdir
- `-w $(pwd)` set container workdir to host workdir
- `-p 8080:8080` expose container 8080 via docker-machine 8080
- `devopspec/engine` use latest [devopspec/engine](https://hub.docker.com/r/devopspec/engine/) image

### 2) Use the engine ReST API 

list dev ops
```
curl $(docker-machine ip):8080/dev-ops
```

run the dev op named `unit-test`
```
curl -X POST $(docker-machine ip):8080/dev-ops/unit-test/runs
```

list pipelines
```
curl $(docker-machine ip):8080/pipelines
```

run the pipeline named `all-tests`
```
curl -X POST $(docker-machine ip):8080/pipelines/all-tests/runs
```

and so on...


# Contributing

refer to [CONTRIBUTING.md](CONTRIBUTING.md)

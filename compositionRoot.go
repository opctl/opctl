package main

import (
  "github.com/opspec-io/engine/core"
  "github.com/opspec-io/engine/tcp"
  dockerPkg "github.com/opspec-io/engine/docker"
)

type compositionRoot interface {
  TcpApi() tcp.Api
}

func newCompositionRoot(
) (compositionRoot compositionRoot, err error) {

  containerEngine, err := dockerPkg.New()
  if (nil != err) {
    return
  }

  compositionRoot = &_compositionRoot{
    tcpApi:tcp.New(
      core.New(
        containerEngine,
      ),
    ),
  }

  return

}

type _compositionRoot struct {
  tcpApi tcp.Api
}

func (_compositionRoot _compositionRoot) TcpApi() tcp.Api {
  return _compositionRoot.tcpApi
}

package main

import (
  "github.com/opspec-io/engine/core"
  "github.com/opspec-io/engine/tcp"
  dockerComposeContainerEngine "github.com/opspec-io/engine/core/adapters/containerengine/dockercompose"
)

type compositionRoot interface {
  TcpApi() tcp.Api
}

func newCompositionRoot(
) (compositionRoot compositionRoot, err error) {

  containerEngine, err := dockerComposeContainerEngine.New()
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

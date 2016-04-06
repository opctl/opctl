package main

import (
  "github.com/dev-op-spec/engine/core"
  "github.com/dev-op-spec/engine/tcp"
  dockerComposeContainerEngine "github.com/dev-op-spec/engine/core/adapters/containerengine/dockercompose"
  osFilesys "github.com/dev-op-spec/engine/core/adapters/filesys/os"
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

  filesys := osFilesys.New()

  coreApi, err := core.New(
    containerEngine,
    filesys,
  )
  if (nil != err) {
    return
  }

  compositionRoot = &_compositionRoot{
    tcpApi:tcp.New(coreApi),
  }

  return

}

type _compositionRoot struct {
  tcpApi tcp.Api
}

func (_compositionRoot _compositionRoot) TcpApi() tcp.Api {
  return _compositionRoot.tcpApi
}

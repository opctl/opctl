package main

import (
  "github.com/open-devops/engine/core"
  "github.com/open-devops/engine/tcp"
  dockerComposeContainerEngine "github.com/open-devops/engine/core/adapters/containerengine/dockercompose"
  osFilesys "github.com/open-devops/engine/core/adapters/filesys/os"
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

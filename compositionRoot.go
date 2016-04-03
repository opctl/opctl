package main

import (
  "github.com/dev-op-spec/engine/core"
  "github.com/dev-op-spec/engine/rest"
  dockerComposeContainerEngine "github.com/dev-op-spec/engine/core/adapters/containerengine/dockercompose"
  osFilesys "github.com/dev-op-spec/engine/core/adapters/filesys/os"
)

type compositionRoot interface {
  RestApi() rest.Api
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
    restApi:rest.New(coreApi),
  }

  return

}

type _compositionRoot struct {
  restApi rest.Api
}

func (_compositionRoot _compositionRoot) RestApi() rest.Api {
  return _compositionRoot.restApi
}

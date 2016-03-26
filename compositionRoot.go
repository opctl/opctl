package main

import (
  "github.com/dev-op-spec/engine/core"
  "github.com/dev-op-spec/engine/rest"
  "github.com/dev-op-spec/engine/core/adapters/containerengine/dockercompose"
  "github.com/dev-op-spec/engine/core/adapters/filesys/os"
)

type compositionRoot interface {
  RestApi() rest.Api
}

func newCompositionRoot(
) (compositionRoot compositionRoot, err error) {

  containerEngine, err := dockercompose.New()
  if (nil != err) {
    return
  }

  filesys := os.New()

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

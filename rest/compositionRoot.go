package rest

import (
  "github.com/dev-op-spec/engine/core"
)

type compositionRoot interface {
  CoreApi() core.Api
}

func newCompositionRoot(
coreApi core.Api,
) (compositionRoot compositionRoot) {

  compositionRoot = &_compositionRoot{
    coreApi:coreApi,
  }

  return

}

type _compositionRoot struct {
  coreApi core.Api
}

func (_compositionRoot _compositionRoot) CoreApi() core.Api {
  return _compositionRoot.coreApi
}

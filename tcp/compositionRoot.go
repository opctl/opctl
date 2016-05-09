package tcp

import (
  "github.com/dev-op-spec/engine/core"
  "net/http"
)

type compositionRoot interface {
  AddOpHandler() http.Handler
  AddSubOpHandler() http.Handler
  GetLivenessHandler() http.Handler
  GetEventStreamHandler() http.Handler
  ListOpsHandler() http.Handler
  RunOpHandler() http.Handler
  SetDescriptionOfOpHandler() http.Handler
}

func newCompositionRoot(
coreApi core.Api,
) (compositionRoot compositionRoot) {

  compositionRoot = &_compositionRoot{
    addOpHandler:newAddOpHandler(coreApi),
    addSubOpHandler:newAddSubOpHandler(coreApi),
    getLivenessHandler:newGetLivenessHandler(),
    getEventStreamHandler:newGetEventStreamHandler(coreApi),
    listOpsHandler:newListOpsHandler(coreApi),
    runOpHandler:newRunOpHandler(coreApi),
    setDescriptionOfOpHandler:newSetDescriptionOfOpHandler(coreApi),
  }

  return

}

type _compositionRoot struct {
  addOpHandler              http.Handler
  addSubOpHandler           http.Handler
  getLivenessHandler          http.Handler
  getEventStreamHandler     http.Handler
  listOpsHandler            http.Handler
  runOpHandler              http.Handler
  setDescriptionOfOpHandler http.Handler
}

func (this _compositionRoot) AddOpHandler(
) http.Handler {
  return this.addOpHandler
}

func (this _compositionRoot) AddSubOpHandler(
) http.Handler {
  return this.addSubOpHandler
}

func (this _compositionRoot) GetLivenessHandler(
) http.Handler {
  return this.getLivenessHandler
}

func (this _compositionRoot) GetEventStreamHandler(
) http.Handler {
  return this.getEventStreamHandler
}

func (this _compositionRoot) ListOpsHandler(
) http.Handler {
  return this.listOpsHandler
}

func (this _compositionRoot) RunOpHandler(
) http.Handler {
  return this.runOpHandler
}

func (this _compositionRoot) SetDescriptionOfOpHandler(
) http.Handler {
  return this.setDescriptionOfOpHandler
}

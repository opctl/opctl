package tcp

import (
  "github.com/opctl/engine/core"
  "net/http"
)

type compositionRoot interface {
  AddOpHandler() http.Handler
  AddSubOpHandler() http.Handler
  GetLivenessHandler() http.Handler
  GetEventStreamHandler() http.Handler
  KillOpRunHandler() http.Handler
  ListOpsHandler() http.Handler
  RunOpHandler() http.Handler
}

func newCompositionRoot(
coreApi core.Api,
) (compositionRoot compositionRoot) {

  compositionRoot = &_compositionRoot{
    addOpHandler:newAddOpHandler(coreApi),
    addSubOpHandler:newAddSubOpHandler(coreApi),
    getLivenessHandler:newGetLivenessHandler(),
    getEventStreamHandler:newGetEventStreamHandler(coreApi),
    killOpRunHandler:newKillOpRunHandler(coreApi),
    listOpsHandler:newListOpsHandler(coreApi),
    runOpHandler:newRunOpHandler(coreApi),
  }

  return

}

type _compositionRoot struct {
  addOpHandler          http.Handler
  addSubOpHandler       http.Handler
  getLivenessHandler    http.Handler
  getEventStreamHandler http.Handler
  killOpRunHandler      http.Handler
  listOpsHandler        http.Handler
  runOpHandler          http.Handler
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

func (this _compositionRoot) KillOpRunHandler(
) http.Handler {
  return this.killOpRunHandler
}

func (this _compositionRoot) ListOpsHandler(
) http.Handler {
  return this.listOpsHandler
}

func (this _compositionRoot) RunOpHandler(
) http.Handler {
  return this.runOpHandler
}

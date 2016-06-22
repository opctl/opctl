package tcp

import (
  "github.com/opctl/engine/core"
  "net/http"
)

type compositionRoot interface {
  GetEventStreamHandler() http.Handler
  GetLivenessHandler() http.Handler
  KillOpRunHandler() http.Handler
  RunOpHandler() http.Handler
}

func newCompositionRoot(
coreApi core.Api,
) (compositionRoot compositionRoot) {

  compositionRoot = &_compositionRoot{
    getEventStreamHandler:newGetEventStreamHandler(coreApi),
    getLivenessHandler:newGetLivenessHandler(),
    killOpRunHandler:newKillOpRunHandler(coreApi),
    runOpHandler:newRunOpHandler(coreApi),
  }

  return

}

type _compositionRoot struct {
  getLivenessHandler    http.Handler
  getEventStreamHandler http.Handler
  killOpRunHandler      http.Handler
  runOpHandler          http.Handler
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

func (this _compositionRoot) RunOpHandler(
) http.Handler {
  return this.runOpHandler
}

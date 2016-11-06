package tcp

import (
  "github.com/opspec-io/engine/core"
  "net/http"
)

type compositionRoot interface {
  GetEventStreamHandler() http.Handler
  GetLivenessHandler() http.Handler
  KillOpRunHandler() http.Handler
  StartOpRunHandler() http.Handler
}

func newCompositionRoot(
coreApi core.Core,
) (compositionRoot compositionRoot) {

  compositionRoot = &_compositionRoot{
    getEventStreamHandler:newGetEventStreamHandler(coreApi),
    getLivenessHandler:newGetLivenessHandler(),
    killOpRunHandler:newKillOpRunHandler(coreApi),
    startOpRunHandler:newStartOpRunHandler(coreApi),
  }

  return

}

type _compositionRoot struct {
  getLivenessHandler    http.Handler
  getEventStreamHandler http.Handler
  killOpRunHandler      http.Handler
  startOpRunHandler          http.Handler
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

func (this _compositionRoot) StartOpRunHandler(
) http.Handler {
  return this.startOpRunHandler
}

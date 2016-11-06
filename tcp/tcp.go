package tcp

import (
  "github.com/opspec-io/engine/core"
  "net/http"
  "github.com/gorilla/mux"
)

type Api interface {
  Start()
}

func New(
coreApi core.Core,
) Api {

  return &_api{
    compositionRoot:newCompositionRoot(coreApi),
  }

}

type _api struct {
  compositionRoot compositionRoot
}

func (this _api) Start(
) {

  router := mux.NewRouter()

  router.Handle(
    getEventStreamRelUrlTemplate,
    this.compositionRoot.GetEventStreamHandler(),
  ).Methods(http.MethodGet)

  router.Handle(
    getLivenessRelUrlTemplate,
    this.compositionRoot.GetLivenessHandler(),
  ).Methods(http.MethodGet)

  router.Handle(
    killOpRunRelUrlTemplate,
    this.compositionRoot.KillOpRunHandler(),
  ).Methods(http.MethodPost)

  router.Handle(
    startOpRunRelUrlTemplate,
    this.compositionRoot.StartOpRunHandler(),
  ).Methods(http.MethodPost)

  http.ListenAndServe(":42224", router)

}

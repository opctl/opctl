package tcp

import (
  "github.com/opctl/engine/core"
  "net/http"
  "github.com/chrisdostert/mux"
)

type Api interface {
  Start()
}

func New(
coreApi core.Api,
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
    addSubOpRelUrlTemplate,
    this.compositionRoot.AddSubOpHandler(),
  ).Methods(http.MethodPost)

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
    runOpRelUrlTemplate,
    this.compositionRoot.RunOpHandler(),
  ).Methods(http.MethodPost)

  router.PathPrefix("/").Handler(http.FileServer(http.Dir("./swagger/")))

  http.ListenAndServe(":42224", router)

}

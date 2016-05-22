package tcp

import (
  "github.com/chrisdostert/mux"
  "github.com/opctl/engine/core"
  "net/http"
  "github.com/codegangsta/negroni"
  "github.com/rs/cors"
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
    addOpRelUrlTemplate,
    this.compositionRoot.AddOpHandler(),
  ).Methods(http.MethodPost)

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
    listOpsRelUrlTemplate,
    this.compositionRoot.ListOpsHandler(),
  ).Methods(http.MethodGet)

  router.Handle(
    runOpRelUrlTemplate,
    this.compositionRoot.RunOpHandler(),
  ).Methods(http.MethodPost)

  router.Handle(
    setDescriptionOfOpRelUrlTemplate,
    this.compositionRoot.SetDescriptionOfOpHandler(),
  ).Methods(http.MethodPut)

  n := negroni.Classic()

  n.Use(cors.New(cors.Options{
    AllowedOrigins: []string{"*"},
    AllowedMethods:[]string{
      http.MethodGet,
      http.MethodPost,
      http.MethodPut,
      http.MethodOptions,
    },
    AllowedHeaders:[]string{"*"},
  }))

  n.Use(negroni.NewStatic(http.Dir("swagger")))

  n.UseHandler(router)

  n.Run(":42224")

}

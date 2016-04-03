package rest

import (
  "github.com/chrisdostert/mux"
  "github.com/dev-op-spec/engine/core"
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
    listOperationsRelUrlTemplate,
    this.compositionRoot.ListOperationsHandler(),
  ).Methods(http.MethodGet)

  router.Handle(
    addOperationRelUrlTemplate,
    this.compositionRoot.AddOperationHandler(),
  ).Methods(http.MethodPost)

  router.Handle(
    setDescriptionOfOperationRelUrlTemplate,
    this.compositionRoot.SetDescriptionOfOperationHandler(),
  ).Methods(http.MethodPut)

  router.Handle(
    addSubOperationRelUrlTemplate,
    this.compositionRoot.AddSubOperationHandler(),
  ).Methods(http.MethodPost)

  router.Handle(
    runOperationRelUrlTemplate,
    this.compositionRoot.RunOperationHandler(),
  ).Methods(http.MethodPost)

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

  n.Run(":8080")

}

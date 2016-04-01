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
    listDevOpsRelUrlTemplate,
    this.compositionRoot.ListDevOpsHandler(),
  ).Methods("GET")

  router.Handle(
    addDevOpRelUrlTemplate,
    this.compositionRoot.AddDevOpHandler(),
  ).Methods("POST")

  router.Handle(
    setDescriptionOfDevOpRelUrlTemplate,
    this.compositionRoot.SetDescriptionOfDevOpHandler(),
  ).Methods("PUT")

  router.Handle(
    runDevOpRelUrlTemplate,
    this.compositionRoot.RunDevOpHandler(),
  ).Methods("POST")

  router.Handle(
    listPipelinesRelUrlTemplate,
    this.compositionRoot.ListPipelinesHandler(),
  ).Methods("GET")

  router.Handle(
    addPipelineRelUrlTemplate,
    this.compositionRoot.AddPipelineHandler(),
  ).Methods("POST")

  router.Handle(
    setDescriptionOfPipelineRelUrlTemplate,
    this.compositionRoot.SetDescriptionOfPipelineHandler(),
  ).Methods("PUT")

  router.Handle(
    addStageToPipelineRelUrlTemplate,
    this.compositionRoot.AddStageToPipelineHandler(),
  ).Methods("POST")

  router.Handle(
    runPipelineRelUrlTemplate,
    this.compositionRoot.RunPipelineHandler(),
  ).Methods("POST")

  n := negroni.Classic()

  n.Use(cors.New(cors.Options{
    AllowedOrigins: []string{"*"},
    AllowedMethods:[]string{"GET", "POST", "PUT", "OPTIONS"},
    AllowedHeaders:[]string{"*"},
  }))

  n.Use(negroni.NewStatic(http.Dir("swagger")))

  n.UseHandler(router)

  n.Run(":8080")

}

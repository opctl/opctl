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
    "/project/{projectUrl}/dev-ops",
    this.compositionRoot.ListDevOpsHandler(),
  ).Methods("GET")

  router.Handle(
    "/project/{projectUrl}/dev-ops",
    this.compositionRoot.AddDevOpHandler(),
  ).Methods("POST")

  router.Handle(
    "/project/{projectUrl}/dev-ops/{devOpName}",
    this.compositionRoot.SetDescriptionOfDevOpHandler(),
  ).Methods("PUT")

  router.Handle(
    "/project/{projectUrl}/dev-ops/{devOpName}/runs",
    this.compositionRoot.RunDevOpHandler(),
  ).Methods("POST")

  router.Handle(
    "/project/{projectUrl}/pipelines",
    this.compositionRoot.ListPipelinesHandler(),
  ).Methods("GET")

  router.Handle(
    "/project/{projectUrl}/pipelines",
    this.compositionRoot.AddPipelineHandler(),
  ).Methods("POST")

  router.Handle(
    "/project/{projectUrl}/pipelines/{pipelineName}",
    this.compositionRoot.SetDescriptionOfPipelineHandler(),
  ).Methods("PUT")

  router.Handle(
    "/project/{projectUrl}/pipelines/{pipelineName}",
    this.compositionRoot.AddStageToPipelineHandler(),
  ).Methods("POST")

  router.Handle(
    "/project/{projectUrl}/pipelines/{pipelineName}/runs",
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

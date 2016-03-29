package rest

import (
  "github.com/chrisdostert/mux"
  "github.com/dev-op-spec/engine/core"
  "net/http"
  "fmt"
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

func MyHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  fmt.Fprintf(w, "Hello %v\n", vars["id"])
}

func (this _api) Start(
) {

  router := mux.NewRouter()

  router.HandleFunc("/some/page/{id}", MyHandler).Methods("GET")

  router.Handle(
    "/project/{projectRootUrl}/dev-ops",
    this.compositionRoot.ListDevOpsHandler(),
  ).Methods("GET")

  router.Handle(
    "/project/{projectRootUrl}/dev-ops",
    this.compositionRoot.AddDevOpHandler(),
  ).Methods("POST")

  router.Handle(
    "/project/{projectRootUrl}/dev-ops/{devOpName}",
    this.compositionRoot.SetDescriptionOfDevOpHandler(),
  ).Methods("PUT")

  router.Handle(
    "/project/{projectRootUrl}/dev-ops/{devOpName}/runs",
    this.compositionRoot.RunDevOpHandler(),
  ).Methods("POST")

  router.Handle(
    "/project/{projectRootUrl}/pipelines",
    this.compositionRoot.ListPipelinesHandler(),
  ).Methods("GET")

  router.Handle(
    "/project/{projectRootUrl}/pipelines",
    this.compositionRoot.AddPipelineHandler(),
  ).Methods("POST")

  router.Handle(
    "/project/{projectRootUrl}/pipelines/{pipelineName}",
    this.compositionRoot.SetDescriptionOfPipelineHandler(),
  ).Methods("PUT")

  router.Handle(
    "/project/{projectRootUrl}/pipelines/{pipelineName}",
    this.compositionRoot.AddStageToPipelineHandler(),
  ).Methods("POST")

  router.Handle(
    "/project/{projectRootUrl}/pipelines/{pipelineName}/runs",
    this.compositionRoot.RunPipelineHandler(),
  ).Methods("POST")

  n := negroni.Classic()

  n.Use(cors.New(cors.Options{
    AllowedOrigins: []string{"*"},
    AllowedMethods:[]string{"GET", "POST", "PUT", "OPTIONS"},
    AllowedHeaders:[]string{"*"},
    Debug:true,
  }))

  n.Use(negroni.NewStatic(http.Dir("swagger")))

  n.UseHandler(router)

  n.Run(":8080")

}

package rest

import (
  "net/http"
  "net/url"
  "github.com/chrisdostert/mux"
  "github.com/dev-op-spec/engine/core"
  "encoding/json"
  "github.com/dev-op-spec/engine/core/models"
)

func newAddPipelineHandler(
coreApi core.Api,
) http.Handler {

  return &addPipelineHandler{
    coreApi:coreApi,
  }

}

type addPipelineHandler struct {
  coreApi core.Api
}

func (this addPipelineHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

  addPipelineReq := models.AddPipelineReq{}

  err := json.NewDecoder(r.Body).Decode(&addPipelineReq)
  if (nil != err) {
    panic(err)
  }

  addPipelineReq.PathToProjectRootDir, err = url.QueryUnescape(mux.Vars(r)["projectRootUrl"])
  if (nil != err) {
    panic(err)
  }

  err = this.coreApi.AddPipeline(addPipelineReq)
  if (nil != err) {
    panic(err)
  }

}

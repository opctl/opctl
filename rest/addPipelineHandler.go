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
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

  var unEscapedProjectUrl string
  unEscapedProjectUrl, err = url.QueryUnescape(mux.Vars(r)["projectUrl"])
  if (nil != err) {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

  addPipelineReq.ProjectUrl, err = models.NewProjectUrl(unEscapedProjectUrl)
  if (nil != err) {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

  err = this.coreApi.AddPipeline(addPipelineReq)
  if (nil != err) {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

}

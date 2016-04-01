package rest

import (
  "net/http"
  "net/url"
  "github.com/chrisdostert/mux"
  "github.com/dev-op-spec/engine/core"
  "encoding/json"
  "github.com/dev-op-spec/engine/core/models"
)

func newRunPipelineHandler(
coreApi core.Api,
) http.Handler {

  return &runPipelineHandler{
    coreApi:coreApi,
  }

}

type runPipelineHandler struct {
  coreApi core.Api
}

func (this runPipelineHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

  runPipelineReq := models.RunPipelineReq{}

  unEscapedProjectUrl, err := url.QueryUnescape(mux.Vars(r)["projectUrl"])
  if (nil != err) {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

  runPipelineReq.ProjectUrl, err = models.NewProjectUrl(unEscapedProjectUrl)
  if (nil != err) {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

  runPipelineReq.PipelineName = mux.Vars(r)["pipelineName"]

  var pipelineRun models.PipelineRunView
  pipelineRun, err = this.coreApi.RunPipeline(runPipelineReq)
  if (nil != err) {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json; charset=UTF-8")

  err = json.NewEncoder(w).Encode(pipelineRun)
  if (nil != err) {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

}

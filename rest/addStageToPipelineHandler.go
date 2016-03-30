package rest

import (
  "net/http"
  "net/url"
  "github.com/chrisdostert/mux"
  "github.com/dev-op-spec/engine/core"
  "encoding/json"
  "github.com/dev-op-spec/engine/core/models"
)

func newAddStageToPipelineHandler(
coreApi core.Api,
) http.Handler {

  return &addStageToPipelineHandler{
    coreApi:coreApi,
  }

}

type addStageToPipelineHandler struct {
  coreApi core.Api
}

func (this addStageToPipelineHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

  addStageToPipelineReq := models.AddStageToPipelineReq{}

  err := json.NewDecoder(r.Body).Decode(&addStageToPipelineReq)
  if (nil != err) {
    panic(err)
  }

  var unEscapedProjectUrl string
  unEscapedProjectUrl, err = url.QueryUnescape(mux.Vars(r)["projectUrl"])
  if (nil != err) {
    panic(err)
  }

  addStageToPipelineReq.ProjectUrl, err = models.NewProjectUrl(unEscapedProjectUrl)
  if (nil != err) {
    panic(err)
  }

  addStageToPipelineReq.PipelineName = mux.Vars(r)["pipelineName"]

  err = this.coreApi.AddStageToPipeline(addStageToPipelineReq)
  if (nil != err) {
    panic(err)
  }

}

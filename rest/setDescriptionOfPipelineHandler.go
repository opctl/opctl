package rest

import (
  "net/http"
  "net/url"
  "github.com/chrisdostert/mux"
  "github.com/dev-op-spec/engine/core"
  "encoding/json"
  "github.com/dev-op-spec/engine/core/models"
)

func newSetDescriptionOfPipelineHandler(
coreApi core.Api,
) http.Handler {

  return &setDescriptionOfPipelineHandler{
    coreApi:coreApi,
  }

}

type setDescriptionOfPipelineHandler struct {
  coreApi core.Api
}

func (this setDescriptionOfPipelineHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

  setDescriptionOfPipelineReq := models.SetDescriptionOfPipelineReq{}

  err := json.NewDecoder(r.Body).Decode(&setDescriptionOfPipelineReq)
  if (nil != err) {
    panic(err)
  }

  var unEscapedProjectUrl string
  unEscapedProjectUrl, err = url.QueryUnescape(mux.Vars(r)["projectUrl"])
  if (nil != err) {
    panic(err)
  }

  setDescriptionOfPipelineReq.ProjectUrl, err = models.NewProjectUrl(unEscapedProjectUrl)
  if (nil != err) {
    panic(err)
  }

  setDescriptionOfPipelineReq.PipelineName = mux.Vars(r)["pipelineName"]

  err = this.coreApi.SetDescriptionOfPipeline(setDescriptionOfPipelineReq)
  if (nil != err) {
    panic(err)
  }

}

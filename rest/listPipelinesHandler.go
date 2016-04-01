package rest

import (
  "net/http"
  "net/url"
  "github.com/chrisdostert/mux"
  "github.com/dev-op-spec/engine/core"
  "encoding/json"
  "github.com/dev-op-spec/engine/core/models"
)

func newListPipelinesHandler(
coreApi core.Api,
) http.Handler {

  return &listPipelinesHandler{
    coreApi:coreApi,
  }

}

type listPipelinesHandler struct {
  coreApi core.Api
}

func (this listPipelinesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

  unEscapedProjectUrl, err := url.QueryUnescape(mux.Vars(r)["projectUrl"])
  if (nil != err) {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

  var projectUrl *models.ProjectUrl
  projectUrl, err = models.NewProjectUrl(unEscapedProjectUrl)
  if (nil != err) {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

  pipelines, err := this.coreApi.ListPipelines(projectUrl)
  if (nil != err) {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json; charset=UTF-8")

  err = json.NewEncoder(w).Encode(pipelines)
  if (nil != err) {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

}

package rest

import (
  "net/http"
  "net/url"
  "github.com/chrisdostert/mux"
  "github.com/dev-op-spec/engine/core"
  "encoding/json"
  "github.com/dev-op-spec/engine/core/models"
)

func newListDevOpsHandler(
coreApi core.Api,
) http.Handler {

  return &listDevOpsHandler{
    coreApi:coreApi,
  }

}

type listDevOpsHandler struct {
  coreApi core.Api
}

func (this listDevOpsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

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

  devOps, err := this.coreApi.ListDevOps(projectUrl)
  if (nil != err) {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json; charset=UTF-8")

  err = json.NewEncoder(w).Encode(devOps)
  if (nil != err) {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

}

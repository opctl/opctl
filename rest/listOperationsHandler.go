package rest

import (
  "net/http"
  "net/url"
  "github.com/chrisdostert/mux"
  "github.com/dev-op-spec/engine/core"
  "encoding/json"
  "github.com/dev-op-spec/engine/core/models"
)

func newListOperationsHandler(
coreApi core.Api,
) http.Handler {

  return &listOperationsHandler{
    coreApi:coreApi,
  }

}

type listOperationsHandler struct {
  coreApi core.Api
}

func (this listOperationsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

  unEscapedProjectUrl, err := url.QueryUnescape(mux.Vars(r)["projectUrl"])
  if (nil != err) {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

  var projectUrl *models.Url
  projectUrl, err = models.NewUrl(unEscapedProjectUrl)
  if (nil != err) {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

  operations, err := this.coreApi.ListOperations(projectUrl)
  if (nil != err) {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json; charset=UTF-8")

  err = json.NewEncoder(w).Encode(operations)
  if (nil != err) {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

}

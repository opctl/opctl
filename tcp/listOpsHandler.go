package tcp

import (
  "net/http"
  "net/url"
  "github.com/chrisdostert/mux"
  "github.com/opctl/engine/core"
  "encoding/json"
  "github.com/opctl/engine/core/models"
)

func newListOpsHandler(
coreApi core.Api,
) http.Handler {

  return &listOpsHandler{
    coreApi:coreApi,
  }

}

type listOpsHandler struct {
  coreApi core.Api
}

func (this listOpsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

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

  ops, err := this.coreApi.ListOps(projectUrl)
  if (nil != err) {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json; charset=UTF-8")

  err = json.NewEncoder(w).Encode(ops)
  if (nil != err) {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

}

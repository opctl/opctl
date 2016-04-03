package rest

import (
  "net/http"
  "net/url"
  "github.com/chrisdostert/mux"
  "github.com/dev-op-spec/engine/core"
  "encoding/json"
  "github.com/dev-op-spec/engine/core/models"
)

func newAddOperationHandler(
coreApi core.Api,
) http.Handler {

  return &addOperationHandler{
    coreApi:coreApi,
  }

}

type addOperationHandler struct {
  coreApi core.Api
}

func (this addOperationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

  addOperationReq := models.AddOperationReq{}

  err := json.NewDecoder(r.Body).Decode(&addOperationReq)
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

  addOperationReq.ProjectUrl, err = models.NewProjectUrl(unEscapedProjectUrl)
  if (nil != err) {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

  err = this.coreApi.AddOperation(addOperationReq)
  if (nil != err) {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

}

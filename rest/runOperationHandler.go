package rest

import (
  "net/http"
  "net/url"
  "github.com/chrisdostert/mux"
  "github.com/dev-op-spec/engine/core"
  "encoding/json"
  "github.com/dev-op-spec/engine/core/models"
)

func newRunOperationHandler(
coreApi core.Api,
) http.Handler {

  return &runOperationHandler{
    coreApi:coreApi,
  }

}

type runOperationHandler struct {
  coreApi core.Api
}

func (this runOperationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

  runOperationReq := models.RunOperationReq{}

  unEscapedProjectUrl, err := url.QueryUnescape(mux.Vars(r)["projectUrl"])
  if (nil != err) {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

  runOperationReq.ProjectUrl, err = models.NewProjectUrl(unEscapedProjectUrl)
  if (nil != err) {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

  runOperationReq.OperationName = mux.Vars(r)["operationName"]

  var operationRun models.OperationRunDetailedView
  operationRun, err = this.coreApi.RunOperation(runOperationReq)
  if (nil != err) {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json; charset=UTF-8")

  err = json.NewEncoder(w).Encode(operationRun)
  if (nil != err) {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

}

package rest

import (
  "net/http"
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

  err := json.NewDecoder(r.Body).Decode(&runOperationReq)
  if (nil != err) {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

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

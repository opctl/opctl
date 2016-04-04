package rest

import (
  "net/http"
  "net/url"
  "github.com/chrisdostert/mux"
  "github.com/dev-op-spec/engine/core"
  "encoding/json"
  "github.com/dev-op-spec/engine/core/models"
)

func newSetDescriptionOfOperationHandler(
coreApi core.Api,
) http.Handler {

  return &setDescriptionOfOperationHandler{
    coreApi:coreApi,
  }

}

type setDescriptionOfOperationHandler struct {
  coreApi core.Api
}

func (this setDescriptionOfOperationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

  setDescriptionOfOperationReq := models.SetDescriptionOfOperationReq{}

  err := json.NewDecoder(r.Body).Decode(&setDescriptionOfOperationReq)
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

  setDescriptionOfOperationReq.ProjectUrl, err = models.NewUrl(unEscapedProjectUrl)
  if (nil != err) {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

  operationName := mux.Vars(r)["operationName"]
  setDescriptionOfOperationReq.OperationName = &operationName

  err = this.coreApi.SetDescriptionOfOperation(setDescriptionOfOperationReq)
  if (nil != err) {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

}

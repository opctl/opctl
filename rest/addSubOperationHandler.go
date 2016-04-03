package rest

import (
  "net/http"
  "net/url"
  "github.com/chrisdostert/mux"
  "github.com/dev-op-spec/engine/core"
  "encoding/json"
  "github.com/dev-op-spec/engine/core/models"
)

func newAddSubOperationHandler(
coreApi core.Api,
) http.Handler {

  return &addSubOperationHandler{
    coreApi:coreApi,
  }

}

type addSubOperationHandler struct {
  coreApi core.Api
}

func (this addSubOperationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

  addSubOperationReq := models.AddSubOperationReq{}

  err := json.NewDecoder(r.Body).Decode(&addSubOperationReq)
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

  addSubOperationReq.ProjectUrl, err = models.NewProjectUrl(unEscapedProjectUrl)
  if (nil != err) {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

  addSubOperationReq.OperationName = mux.Vars(r)["operationName"]

  err = this.coreApi.AddSubOperation(addSubOperationReq)
  if (nil != err) {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

}

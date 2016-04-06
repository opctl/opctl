package tcp

import (
  "net/http"
  "net/url"
  "github.com/chrisdostert/mux"
  "github.com/dev-op-spec/engine/core"
  "encoding/json"
  "github.com/dev-op-spec/engine/core/models"
)

func newSetDescriptionOfOpHandler(
coreApi core.Api,
) http.Handler {

  return &setDescriptionOfOpHandler{
    coreApi:coreApi,
  }

}

type setDescriptionOfOpHandler struct {
  coreApi core.Api
}

func (this setDescriptionOfOpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

  setDescriptionOfOpReq := models.SetDescriptionOfOpReq{}

  err := json.NewDecoder(r.Body).Decode(&setDescriptionOfOpReq)
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

  setDescriptionOfOpReq.ProjectUrl, err = models.NewUrl(unEscapedProjectUrl)
  if (nil != err) {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

  opName := mux.Vars(r)["opName"]
  setDescriptionOfOpReq.OpName = &opName

  err = this.coreApi.SetDescriptionOfOp(setDescriptionOfOpReq)
  if (nil != err) {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

}

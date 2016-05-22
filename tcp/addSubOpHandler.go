package tcp

import (
  "net/http"
  "net/url"
  "github.com/chrisdostert/mux"
  "github.com/opctl/engine/core"
  "encoding/json"
  "github.com/opctl/engine/core/models"
)

func newAddSubOpHandler(
coreApi core.Api,
) http.Handler {

  return &addSubOpHandler{
    coreApi:coreApi,
  }

}

type addSubOpHandler struct {
  coreApi core.Api
}

func (this addSubOpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

  addSubOpReq := models.AddSubOpReq{}

  err := json.NewDecoder(r.Body).Decode(&addSubOpReq)
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

  addSubOpReq.ProjectUrl, err = models.NewUrl(unEscapedProjectUrl)
  if (nil != err) {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

  addSubOpReq.OpName = mux.Vars(r)["opName"]

  err = this.coreApi.AddSubOp(addSubOpReq)
  if (nil != err) {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  w.WriteHeader(http.StatusCreated)

}

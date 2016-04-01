package rest

import (
  "net/http"
  "net/url"
  "github.com/chrisdostert/mux"
  "github.com/dev-op-spec/engine/core"
  "encoding/json"
  "github.com/dev-op-spec/engine/core/models"
)

func newRunDevOpHandler(
coreApi core.Api,
) http.Handler {

  return &runDevOpHandler{
    coreApi:coreApi,
  }

}

type runDevOpHandler struct {
  coreApi core.Api
}

func (this runDevOpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

  runDevOpReq := models.RunDevOpReq{}

  unEscapedProjectUrl, err := url.QueryUnescape(mux.Vars(r)["projectUrl"])
  if (nil != err) {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

  runDevOpReq.ProjectUrl, err = models.NewProjectUrl(unEscapedProjectUrl)
  if (nil != err) {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

  runDevOpReq.DevOpName = mux.Vars(r)["devOpName"]

  var devOpRun models.DevOpRunView
  devOpRun, err = this.coreApi.RunDevOp(runDevOpReq)
  if (nil != err) {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json; charset=UTF-8")

  err = json.NewEncoder(w).Encode(devOpRun)
  if (nil != err) {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

}

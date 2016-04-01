package rest

import (
  "net/http"
  "net/url"
  "github.com/chrisdostert/mux"
  "github.com/dev-op-spec/engine/core"
  "encoding/json"
  "github.com/dev-op-spec/engine/core/models"
)

func newAddDevOpHandler(
coreApi core.Api,
) http.Handler {

  return &addDevOpHandler{
    coreApi:coreApi,
  }

}

type addDevOpHandler struct {
  coreApi core.Api
}

func (this addDevOpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

  addDevOpReq := models.AddDevOpReq{}

  err := json.NewDecoder(r.Body).Decode(&addDevOpReq)
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

  addDevOpReq.ProjectUrl, err = models.NewProjectUrl(unEscapedProjectUrl)
  if (nil != err) {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

  err = this.coreApi.AddDevOp(addDevOpReq)
  if (nil != err) {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

}

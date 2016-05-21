package tcp

import (
  "net/http"
  "net/url"
  "github.com/chrisdostert/mux"
  "github.com/open-devops/engine/core"
  "encoding/json"
  "github.com/open-devops/engine/core/models"
)

func newAddOpHandler(
coreApi core.Api,
) http.Handler {

  return &addOpHandler{
    coreApi:coreApi,
  }

}

type addOpHandler struct {
  coreApi core.Api
}

func (this addOpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

  addOpReq := models.AddOpReq{}

  err := json.NewDecoder(r.Body).Decode(&addOpReq)
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

  addOpReq.ProjectUrl, err = models.NewUrl(unEscapedProjectUrl)
  if (nil != err) {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

  err = this.coreApi.AddOp(addOpReq)
  if (nil != err) {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  w.WriteHeader(http.StatusCreated)

}

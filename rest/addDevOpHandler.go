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
    panic(err)
  }

  addDevOpReq.PathToProjectRootDir, err = url.QueryUnescape(mux.Vars(r)["projectRootUrl"])
  if (nil != err) {
    panic(err)
  }

  err = this.coreApi.AddDevOp(addDevOpReq)
  if (nil != err) {
    panic(err)
  }

}

package rest

import (
  "net/http"
  "net/url"
  "github.com/chrisdostert/mux"
  "github.com/dev-op-spec/engine/core"
  "encoding/json"
)

func newListDevOpsHandler(
coreApi core.Api,
) http.Handler {

  return &listDevOpsHandler{
    coreApi:coreApi,
  }

}

type listDevOpsHandler struct {
  coreApi core.Api
}

func (this listDevOpsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

  projectRootUrl, err := url.QueryUnescape(mux.Vars(r)["projectRootUrl"])
  if (nil != err) {
    panic(err)
  }

  devOps, err := this.coreApi.ListDevOps(projectRootUrl)
  if (nil != err) {
    panic(err)
  }

  w.Header().Set("Content-Type", "application/json; charset=UTF-8")

  err = json.NewEncoder(w).Encode(devOps)
  if (nil != err) {
    panic(err)
  }

}

package tcp

import (
  "net/http"
  "github.com/opspec-io/engine/core"
  "encoding/json"
  "github.com/opspec-io/sdk-golang/pkg/model"
)

func newKillOpRunHandler(
core core.Core,
) http.Handler {

  return &killOpRunHandler{
    core:core,
  }

}

type killOpRunHandler struct {
  core core.Core
}

func (this killOpRunHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

  killOpRunReq := model.KillOpRunReq{}

  err := json.NewDecoder(r.Body).Decode(&killOpRunReq)
  if (nil != err) {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

  this.core.KillOpRun(killOpRunReq)

  w.WriteHeader(http.StatusCreated)
  w.Header().Set("Content-Type", "text/plain; charset=UTF-8")

}

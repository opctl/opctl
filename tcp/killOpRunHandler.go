package tcp

import (
  "net/http"
  "github.com/opspec-io/engine/core"
  "encoding/json"
  "github.com/opspec-io/sdk-golang/pkg/models"
)

func newKillOpRunHandler(
coreApi core.Core,
) http.Handler {

  return &killOpRunHandler{
    coreApi:coreApi,
  }

}

type killOpRunHandler struct {
  coreApi core.Core
}

func (this killOpRunHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

  killOpRunReq := models.KillOpRunReq{}

  err := json.NewDecoder(r.Body).Decode(&killOpRunReq)
  if (nil != err) {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

  this.coreApi.KillOpRun(killOpRunReq)

  w.WriteHeader(http.StatusCreated)
  w.Header().Set("Content-Type", "text/plain; charset=UTF-8")

}

package tcp

import (
  "net/http"
  "github.com/opctl/engine/core"
  "encoding/json"
  "github.com/opctl/engine/core/models"
)

func newKillOpRunHandler(
coreApi core.Api,
) http.Handler {

  return &killOpRunHandler{
    coreApi:coreApi,
  }

}

type killOpRunHandler struct {
  coreApi core.Api
}

func (this killOpRunHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

  killOpRunReq := models.KillOpRunReq{}

  err := json.NewDecoder(r.Body).Decode(&killOpRunReq)
  if (nil != err) {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

  correlationId, err := this.coreApi.KillOpRun(killOpRunReq)
  if (nil != err) {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  w.Header().Set("Correlation-Id", correlationId)

  w.WriteHeader(http.StatusCreated)
  w.Header().Set("Content-Type", "text/plain; charset=UTF-8")

}

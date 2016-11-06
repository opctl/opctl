package tcp

import (
  "net/http"
  "github.com/opspec-io/engine/core"
  "encoding/json"
  "github.com/opspec-io/sdk-golang/pkg/models"
)

func newStartOpRunHandler(
coreApi core.Core,
) http.Handler {

  return &startOpRunHandler{
    coreApi:coreApi,
  }

}

type startOpRunHandler struct {
  coreApi core.Core
}

func (this startOpRunHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

  startOpRunReq := models.StartOpRunReq{}

  err := json.NewDecoder(r.Body).Decode(&startOpRunReq)
  if (nil != err) {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

  opRunId, err := this.coreApi.StartOpRun(startOpRunReq)
  if (nil != err) {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  w.WriteHeader(http.StatusCreated)
  w.Header().Set("Content-Type", "text/plain; charset=UTF-8")

  w.Write([]byte(opRunId))
  if (nil != err) {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

}

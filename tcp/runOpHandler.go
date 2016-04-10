package tcp

import (
  "net/http"
  "github.com/dev-op-spec/engine/core"
  "encoding/json"
  "github.com/dev-op-spec/engine/core/models"
)

func newRunOpHandler(
coreApi core.Api,
) http.Handler {

  return &runOpHandler{
    coreApi:coreApi,
  }

}

type runOpHandler struct {
  coreApi core.Api
}

func (this runOpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

  runOpReq := models.RunOpReq{}

  err := json.NewDecoder(r.Body).Decode(&runOpReq)
  if (nil != err) {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

  opRunId, err := this.coreApi.RunOp(runOpReq)
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

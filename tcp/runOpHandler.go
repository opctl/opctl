package tcp

import (
  "net/http"
  "github.com/opspec-io/engine/core"
  "encoding/json"
  "github.com/opspec-io/sdk-golang/pkg/model"
)

func newStartOpRunHandler(
core core.Core,
) http.Handler {

  return &startOpRunHandler{
    core:core,
  }

}

type startOpRunHandler struct {
  core core.Core
}

func (this startOpRunHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

  startOpRunReq := model.StartOpRunReq{}

  err := json.NewDecoder(r.Body).Decode(&startOpRunReq)
  if (nil != err) {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

  opRunId, err := this.core.StartOpRun(startOpRunReq)
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

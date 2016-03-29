package rest

import (
  "net/http"
  "net/url"
  "github.com/chrisdostert/mux"
  "github.com/dev-op-spec/engine/core"
  "encoding/json"
)

func newListPipelinesHandler(
coreApi core.Api,
) http.Handler {

  return &listPipelinesHandler{
    coreApi:coreApi,
  }

}

type listPipelinesHandler struct {
  coreApi core.Api
}

func (this listPipelinesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

  projectRootUrl, err := url.QueryUnescape(mux.Vars(r)["projectRootUrl"])
  if (nil != err) {
    panic(err)
  }

  pipelines, err := this.coreApi.ListPipelines(projectRootUrl)
  if (nil != err) {
    panic(err)
  }

  w.Header().Set("Content-Type", "application/json; charset=UTF-8")

  err = json.NewEncoder(w).Encode(pipelines)
  if (nil != err) {
    panic(err)
  }

}

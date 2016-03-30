package rest

import (
  "net/http"
  "net/url"
  "github.com/chrisdostert/mux"
  "github.com/dev-op-spec/engine/core"
  "encoding/json"
  "github.com/dev-op-spec/engine/core/models"
)

func newSetDescriptionOfDevOpHandler(
coreApi core.Api,
) http.Handler {

  return &setDescriptionOfDevOpHandler{
    coreApi:coreApi,
  }

}

type setDescriptionOfDevOpHandler struct {
  coreApi core.Api
}

func (this setDescriptionOfDevOpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

  setDescriptionOfDevOpReq := models.SetDescriptionOfDevOpReq{}

  err := json.NewDecoder(r.Body).Decode(&setDescriptionOfDevOpReq)
  if (nil != err) {
    panic(err)
  }

  var unEscapedProjectUrl string
  unEscapedProjectUrl, err = url.QueryUnescape(mux.Vars(r)["projectUrl"])
  if (nil != err) {
    panic(err)
  }

  setDescriptionOfDevOpReq.ProjectUrl, err = models.NewProjectUrl(unEscapedProjectUrl)
  if (nil != err) {
    panic(err)
  }

  setDescriptionOfDevOpReq.DevOpName = mux.Vars(r)["devOpName"]

  err = this.coreApi.SetDescriptionOfDevOp(setDescriptionOfDevOpReq)
  if (nil != err) {
    panic(err)
  }

}

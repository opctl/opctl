package handler

import (
	"encoding/json"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/node/core"
	"net/http"
)

func newKillOpHandler(
	core core.Core,
) http.Handler {
	return killOpHandler{
		core: core,
	}
}

type killOpHandler struct {
	core core.Core
}

func (koh killOpHandler) ServeHTTP(httpResp http.ResponseWriter, httpReq *http.Request) {

	killOpReq := model.KillOpReq{}

	err := json.NewDecoder(httpReq.Body).Decode(&killOpReq)
	if nil != err {
		http.Error(httpResp, err.Error(), http.StatusBadRequest)
		return
	}

	koh.core.KillOp(killOpReq)

	httpResp.WriteHeader(http.StatusCreated)
	httpResp.Header().Set("Content-Type", "text/plain; charset=UTF-8")

}

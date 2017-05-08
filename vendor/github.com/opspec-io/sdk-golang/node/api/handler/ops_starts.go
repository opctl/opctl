package handler

import (
	"encoding/json"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/node/core"
	"net/http"
)

func newStartOpHandler(
	core core.Core,
) http.Handler {
	return startOpHandler{
		core: core,
	}
}

type startOpHandler struct {
	core core.Core
}

func (soh startOpHandler) ServeHTTP(httpResp http.ResponseWriter, httpReq *http.Request) {

	startOpReq := model.StartOpReq{}

	err := json.NewDecoder(httpReq.Body).Decode(&startOpReq)
	if nil != err {
		http.Error(httpResp, err.Error(), http.StatusBadRequest)
		return
	}

	callId, err := soh.core.StartOp(startOpReq)
	if nil != err {
		http.Error(httpResp, err.Error(), http.StatusInternalServerError)
		return
	}

	httpResp.WriteHeader(http.StatusCreated)
	httpResp.Header().Set("Content-Type", "text/plain; charset=UTF-8")

	httpResp.Write([]byte(callId))
	if nil != err {
		http.Error(httpResp, err.Error(), http.StatusInternalServerError)
		return
	}

}

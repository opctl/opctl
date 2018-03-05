package handler

import (
	"encoding/json"
	"github.com/opspec-io/sdk-golang/model"
	"net/http"
)

func (hdlr _handler) ops_starts(
	httpResp http.ResponseWriter,
	httpReq *http.Request,
) {

	startOpReq := model.StartOpReq{}

	err := json.NewDecoder(httpReq.Body).Decode(&startOpReq)
	if nil != err {
		http.Error(httpResp, err.Error(), http.StatusBadRequest)
		return
	}

	callID, err := hdlr.core.StartOp(
		httpReq.Context(),
		startOpReq,
	)
	if nil != err {
		http.Error(httpResp, err.Error(), http.StatusInternalServerError)
		return
	}

	httpResp.WriteHeader(http.StatusCreated)
	httpResp.Header().Set("Content-Type", "text/plain; charset=UTF-8")

	httpResp.Write([]byte(callID))
	if nil != err {
		http.Error(httpResp, err.Error(), http.StatusInternalServerError)
		return
	}

}

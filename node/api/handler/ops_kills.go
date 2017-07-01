package handler

import (
	"encoding/json"
	"github.com/opspec-io/sdk-golang/model"
	"net/http"
)

func (hdlr _handler) ops_kills(
	httpResp http.ResponseWriter,
	httpReq *http.Request,
) {

	killOpReq := model.KillOpReq{}

	err := json.NewDecoder(httpReq.Body).Decode(&killOpReq)
	if nil != err {
		http.Error(httpResp, err.Error(), http.StatusBadRequest)
		return
	}

	hdlr.core.KillOp(killOpReq)

	httpResp.WriteHeader(http.StatusCreated)
	httpResp.Header().Set("Content-Type", "text/plain; charset=UTF-8")

}

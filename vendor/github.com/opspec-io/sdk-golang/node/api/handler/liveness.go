package handler

import (
	"net/http"
)

func (hdlr _handler) liveness(
	httpResp http.ResponseWriter,
	httpReq *http.Request,
) {

	// currently this is a no-op

}

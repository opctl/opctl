package handler

import (
	"net/http"
)

func newGetLivenessHandler() http.Handler {

	return &getLivenessHandler{}

}

type getLivenessHandler struct{}

func (glh getLivenessHandler) ServeHTTP(httpResp http.ResponseWriter, httpReq *http.Request) {

	// currently this is a no-op

}

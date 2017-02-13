package tcp

import (
	"net/http"
)

func newGetLivenessHandler() http.Handler {

	return &getLivenessHandler{}

}

type getLivenessHandler struct{}

func (this getLivenessHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// currently this is a no-op

}

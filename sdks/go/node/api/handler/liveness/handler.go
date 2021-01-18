package liveness

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

import (
	"net/http"

	"github.com/opctl/opctl/sdks/go/internal/urlpath"
	"github.com/opctl/opctl/sdks/go/node"
)

//counterfeiter:generate -o fakes/handler.go . Handler
type Handler interface {
	Handle(
		httpResp http.ResponseWriter,
		httpReq *http.Request,
	)
}

// NewHandler returns an initialized Handler instance
func NewHandler(
	node node.Node,
) Handler {
	return _handler{}
}

type _handler struct {
}

func (hdlr _handler) Handle(
	httpResp http.ResponseWriter,
	httpReq *http.Request,
) {
	pathSegment, err := urlpath.NextSegment(httpReq.URL)
	if nil != err {
		http.Error(httpResp, err.Error(), http.StatusBadRequest)
		return
	}

	if "" != pathSegment {
		http.NotFoundHandler().ServeHTTP(httpResp, httpReq)
		return
	}

	// currently this is a no-op
}

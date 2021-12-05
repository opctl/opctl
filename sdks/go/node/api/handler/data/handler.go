package data

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

import (
	"net/http"

	"github.com/opctl/opctl/sdks/go/internal/urlpath"
	"github.com/opctl/opctl/sdks/go/node"
	"github.com/opctl/opctl/sdks/go/node/api/handler/data/ref"
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
	node node.Core,
) Handler {
	return _handler{
		refHandler: ref.NewHandler(node),
	}
}

type _handler struct {
	refHandler ref.Handler
}

func (hdlr _handler) Handle(
	httpResp http.ResponseWriter,
	httpReq *http.Request,
) {
	pathSegment, err := urlpath.NextSegment(httpReq.URL)
	if err != nil {
		http.Error(httpResp, err.Error(), http.StatusBadRequest)
		return
	}

	switch pathSegment {
	case "":
		http.NotFoundHandler().ServeHTTP(httpResp, httpReq)
		return
	default:
		hdlr.refHandler.Handle(
			pathSegment,
			httpResp,
			httpReq,
		)
	}
}

package ops

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

import (
	"net/http"

	"github.com/opctl/opctl/sdks/go/internal/urlpath"
	"github.com/opctl/opctl/sdks/go/node"
	"github.com/opctl/opctl/sdks/go/node/api/handler/ops/kills"
	"github.com/opctl/opctl/sdks/go/node/api/handler/ops/starts"
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
	return _handler{
		startsHandler: starts.NewHandler(node),
		killsHandler:  kills.NewHandler(node),
	}
}

type _handler struct {
	startsHandler starts.Handler
	killsHandler  kills.Handler
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
	case "kills":
		hdlr.killsHandler.Handle(
			httpResp,
			httpReq,
		)
	case "starts":
		hdlr.startsHandler.Handle(
			httpResp,
			httpReq,
		)
	default:
		http.NotFoundHandler().ServeHTTP(httpResp, httpReq)
		return
	}
}

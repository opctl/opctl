package events

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

import (
	"github.com/opctl/opctl/sdks/go/internal/urlpath"
	"github.com/opctl/opctl/sdks/go/node/api/handler/events/stream"
	"github.com/opctl/opctl/sdks/go/node/core"
	"net/http"
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
	core core.Core,
) Handler {
	return _handler{
		streamHandler: stream.NewHandler(core),
	}
}

type _handler struct {
	streamHandler stream.Handler
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

	switch pathSegment {
	case "stream":
		hdlr.streamHandler.Handle(
			httpResp,
			httpReq,
		)
	default:
		http.NotFoundHandler().ServeHTTP(httpResp, httpReq)
		return
	}
}

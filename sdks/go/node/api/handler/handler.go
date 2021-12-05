package handler

import (
	"net/http"

	"github.com/opctl/opctl/sdks/go/internal/urlpath"
	"github.com/opctl/opctl/sdks/go/node"
	"github.com/opctl/opctl/sdks/go/node/api/handler/auths"
	"github.com/opctl/opctl/sdks/go/node/api/handler/data"
	"github.com/opctl/opctl/sdks/go/node/api/handler/events"
	"github.com/opctl/opctl/sdks/go/node/api/handler/liveness"
	"github.com/opctl/opctl/sdks/go/node/api/handler/ops"
)

var oneMB int64 = 1024 * 1024
var maxReqBytes int64 = 40 * oneMB

// New returns an http server that wraps the given Core op runner with an http
// API. APIClient provides an Node interface for interacting with it.
func New(
	node node.Core,
) http.Handler {
	return _handler{
		authsHandler:    auths.NewHandler(node),
		dataHandler:     data.NewHandler(node),
		eventsHandler:   events.NewHandler(node),
		livenessHandler: liveness.NewHandler(node),
		opsHandler:      ops.NewHandler(node),
	}
}

type _handler struct {
	authsHandler    auths.Handler
	dataHandler     data.Handler
	eventsHandler   events.Handler
	livenessHandler liveness.Handler
	opsHandler      ops.Handler
}

func (hdlr _handler) ServeHTTP(
	httpResp http.ResponseWriter,
	httpReq *http.Request,
) {
	// limit req size to maxReqBytes
	httpReq.Body = http.MaxBytesReader(httpResp, httpReq.Body, maxReqBytes)

	pathSegment, err := urlpath.NextSegment(httpReq.URL)
	if err != nil {
		http.Error(httpResp, err.Error(), http.StatusBadRequest)
		return
	}

	switch pathSegment {
	case "auths":
		hdlr.authsHandler.Handle(httpResp, httpReq)
	case "data":
		hdlr.dataHandler.Handle(httpResp, httpReq)
	case "events":
		hdlr.eventsHandler.Handle(httpResp, httpReq)
	case "liveness":
		hdlr.livenessHandler.Handle(httpResp, httpReq)
	case "ops":
		hdlr.opsHandler.Handle(httpResp, httpReq)
	default:
		http.NotFoundHandler().ServeHTTP(httpResp, httpReq)
	}
}

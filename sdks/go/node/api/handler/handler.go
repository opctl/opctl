package handler

import (
	"net/http"

	"github.com/opctl/opctl/sdks/go/internal/urlpath"
	"github.com/opctl/opctl/sdks/go/node/api/handler/auths"
	"github.com/opctl/opctl/sdks/go/node/api/handler/data"
	"github.com/opctl/opctl/sdks/go/node/api/handler/events"
	"github.com/opctl/opctl/sdks/go/node/api/handler/liveness"
	"github.com/opctl/opctl/sdks/go/node/api/handler/ops"
	"github.com/opctl/opctl/sdks/go/node/api/handler/pkgs"
	"github.com/opctl/opctl/sdks/go/node/core"
)

// New returns an http server that wraps the given Core op runner with an http
// API. APIClient provides an OpNode interface for interacting with it.
func New(
	core core.Core,
) http.Handler {
	return _handler{
		authsHandler:    auths.NewHandler(core),
		dataHandler:     data.NewHandler(core),
		eventsHandler:   events.NewHandler(core),
		livenessHandler: liveness.NewHandler(core),
		opsHandler:      ops.NewHandler(core),
		pkgsHandler:     pkgs.NewHandler(core),
	}
}

type _handler struct {
	authsHandler    auths.Handler
	dataHandler     data.Handler
	eventsHandler   events.Handler
	livenessHandler liveness.Handler
	opsHandler      ops.Handler
	pkgsHandler     pkgs.Handler
}

func (hdlr _handler) ServeHTTP(
	httpResp http.ResponseWriter,
	httpReq *http.Request,
) {
	pathSegment, err := urlpath.NextSegment(httpReq.URL)
	if nil != err {
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
	case "pkgs":
		// deprecated resource
		hdlr.pkgsHandler.Handle(httpResp, httpReq)
	default:
		http.NotFoundHandler().ServeHTTP(httpResp, httpReq)
	}
}

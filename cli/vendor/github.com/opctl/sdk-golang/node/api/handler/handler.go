package handler

import (
	"net/http"

	"github.com/opctl/sdk-golang/node/api/handler/data"
	"github.com/opctl/sdk-golang/node/api/handler/events"
	"github.com/opctl/sdk-golang/node/api/handler/liveness"
	"github.com/opctl/sdk-golang/node/api/handler/ops"
	"github.com/opctl/sdk-golang/node/api/handler/pkgs"
	"github.com/opctl/sdk-golang/node/core"
	"github.com/opctl/sdk-golang/util/urlpath"
)

func New(
	core core.Core,
) http.Handler {
	return _handler{
		dataHandler:     data.NewHandler(core),
		eventsHandler:   events.NewHandler(core),
		livenessHandler: liveness.NewHandler(core),
		opsHandler:      ops.NewHandler(core),
		pkgsHandler:     pkgs.NewHandler(core),
	}
}

type _handler struct {
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

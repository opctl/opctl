package contents

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

import (
	"encoding/json"
	"net/http"

	"github.com/opctl/opctl/sdks/go/internal/urlpath"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node"
	"github.com/opctl/opctl/sdks/go/node/api/handler/pkgs/ref/contents/path"
)

// Handler deprecated
//counterfeiter:generate -o fakes/handler.go . Handler
type Handler interface {
	Handle(
		dataHandle model.DataHandle,
		httpResp http.ResponseWriter,
		httpReq *http.Request,
	)
}

// NewHandler returns an initialized Handler instance
func NewHandler(
	node node.Node,
) Handler {
	return _handler{
		node:        node,
		pathHandler: path.NewHandler(node),
	}
}

type _handler struct {
	node        node.Node
	pathHandler path.Handler
}

func (hdlr _handler) Handle(
	dataHandle model.DataHandle,
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
		dirEntriesList, err := dataHandle.ListDescendants(
			httpReq.Context(),
		)
		if err != nil {
			http.Error(httpResp, err.Error(), http.StatusInternalServerError)
			return
		}

		httpResp.Header().Set("Content-Type", "application/json; charset=UTF-8")

		if err := json.NewEncoder(httpResp).Encode(dirEntriesList); err != nil {
			http.Error(httpResp, err.Error(), http.StatusInternalServerError)
			return
		}
	default:
		hdlr.pathHandler.Handle(
			dataHandle,
			pathSegment,
			httpResp,
			httpReq,
		)
	}
}

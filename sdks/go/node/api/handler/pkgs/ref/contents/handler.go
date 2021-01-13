package contents

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

import (
	"net/http"

	ijson "github.com/golang-interfaces/encoding-ijson"
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
	opNode node.OpNode,
) Handler {
	return _handler{
		opNode:      opNode,
		json:        ijson.New(),
		pathHandler: path.NewHandler(opNode),
	}
}

type _handler struct {
	opNode      node.OpNode
	json        ijson.IJSON
	pathHandler path.Handler
}

func (hdlr _handler) Handle(
	dataHandle model.DataHandle,
	httpResp http.ResponseWriter,
	httpReq *http.Request,
) {
	pathSegment, err := urlpath.NextSegment(httpReq.URL)
	if nil != err {
		http.Error(httpResp, err.Error(), http.StatusBadRequest)
		return
	}

	switch pathSegment {
	case "":
		dirEntriesList, err := dataHandle.ListDescendants(
			httpReq.Context(),
		)
		if nil != err {
			http.Error(httpResp, err.Error(), http.StatusInternalServerError)
			return
		}

		httpResp.Header().Set("Content-Type", "application/json; charset=UTF-8")

		if err := hdlr.json.NewEncoder(httpResp).Encode(dirEntriesList); nil != err {
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

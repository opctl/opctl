package path

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

import (
	"net/http"
	pathPkg "path"
	"time"

	"github.com/golang-interfaces/ihttp"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node"
)

// Handler deprecated
//counterfeiter:generate -o fakes/handler.go . Handler
type Handler interface {
	Handle(
		dataHandle model.DataHandle,
		dataPath string,
		httpResp http.ResponseWriter,
		httpReq *http.Request,
	)
}

// NewHandler returns an initialized Handler instance
func NewHandler(
	node node.Node,
) Handler {
	return _handler{
		node: node,
		http: ihttp.New(),
	}
}

type _handler struct {
	node node.Node
	http ihttp.IHTTP
}

func (hdlr _handler) Handle(
	dataHandle model.DataHandle,
	dataPath string,
	httpResp http.ResponseWriter,
	httpReq *http.Request,
) {
	dirEntryReader, err := dataHandle.GetContent(
		httpReq.Context(),
		dataPath,
	)
	if nil != err {
		http.Error(httpResp, err.Error(), http.StatusInternalServerError)
		return
	}

	hdlr.http.ServeContent(
		httpResp,
		httpReq,
		pathPkg.Base(dataPath),
		time.Time{},
		dirEntryReader,
	)
}

package path

//go:generate counterfeiter -o ./fakeHandler.go --fake-name FakeHandler ./ Handler

import (
	"net/http"
	pathPkg "path"
	"time"

	"github.com/golang-interfaces/ihttp"
	"github.com/opctl/sdk-golang/model"
	"github.com/opctl/sdk-golang/node/core"
)

// Handler deprecated
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
	core core.Core,
) Handler {
	return _handler{
		core: core,
		http: ihttp.New(),
	}
}

type _handler struct {
	core core.Core
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

package ref

//go:generate counterfeiter -o ./fakeHandleGetOrHeader.go --fake-name fakeHandleGetOrHeader ./ handleGetOrHeader

import (
	"net/http"
	"path"
	"time"

	"github.com/golang-interfaces/ihttp"
	"github.com/opctl/sdk-golang/model"
	"github.com/opctl/sdk-golang/node/core"
)

// handleGetOrHeader handles GET or HEAD's
type handleGetOrHeader interface {
	HandleGetOrHead(
		dataHandle model.DataHandle,
		httpResp http.ResponseWriter,
		httpReq *http.Request,
	)
}

// newHandleGetOrHeader returns an initialized handleGetOrHeader instance
func newHandleGetOrHeader(
	core core.Core,
) handleGetOrHeader {
	return _handleGetOrHeader{
		core: core,
		http: ihttp.New(),
	}
}

type _handleGetOrHeader struct {
	core core.Core
	http ihttp.IHTTP
}

func (hg _handleGetOrHeader) HandleGetOrHead(
	dataHandle model.DataHandle,
	httpResp http.ResponseWriter,
	httpReq *http.Request,
) {
	if httpReq.URL.Path != "" {
		http.Error(httpResp, "", http.StatusNotFound)
		return
	}

	if httpReq.Method != http.MethodGet && httpReq.Method != http.MethodHead {
		http.Error(httpResp, "Request MUST be GET or HEAD", http.StatusMethodNotAllowed)
		return
	}

	dirEntryReader, err := dataHandle.GetContent(
		httpReq.Context(),
		"",
	)
	if nil != err {
		http.Error(httpResp, err.Error(), http.StatusInternalServerError)
		return
	}

	hg.http.ServeContent(
		httpResp,
		httpReq,
		path.Base(dataHandle.Ref()),
		time.Time{},
		dirEntryReader,
	)
}

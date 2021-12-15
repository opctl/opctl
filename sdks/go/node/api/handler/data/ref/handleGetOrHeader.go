package ref

import (
	"encoding/json"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/golang-interfaces/ihttp"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node"
)

// handleGetOrHeader handles GET or HEAD's
//counterfeiter:generate -o internal/fakes/handleGetOrHeader.go . handleGetOrHeader
type handleGetOrHeader interface {
	HandleGetOrHead(
		dataHandle model.DataHandle,
		httpResp http.ResponseWriter,
		httpReq *http.Request,
	)
}

// newHandleGetOrHeader returns an initialized handleGetOrHeader instance
func newHandleGetOrHeader(
	node node.Core,
) handleGetOrHeader {
	return _handleGetOrHeader{
		node: node,
		http: ihttp.New(),
	}
}

type _handleGetOrHeader struct {
	node node.Core
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

	dataPath := dataHandle.Path()
	dataFileInfo, err := os.Stat(*dataPath)
	if err != nil {
		http.Error(httpResp, err.Error(), http.StatusInternalServerError)
		return
	}

	if dataFileInfo.IsDir() {
		dirEntriesList, err := dataHandle.ListDescendants(
			httpReq.Context(),
		)
		if err != nil {
			http.Error(httpResp, err.Error(), http.StatusInternalServerError)
			return
		}

		httpResp.Header().Set("Content-Type", "application/vnd.opspec.0.1.6.dir+json; charset=UTF-8")

		if err := json.NewEncoder(httpResp).Encode(dirEntriesList); err != nil {
			http.Error(httpResp, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	dirEntryReader, err := dataHandle.GetContent(
		httpReq.Context(),
		"",
	)
	if err != nil {
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

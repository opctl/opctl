package ref

//go:generate counterfeiter -o ./fakeHandler.go --fake-name FakeHandler ./ Handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/node/core"
	"github.com/opspec-io/sdk-golang/util/urlpath"
)

type Handler interface {
	Handle(
		res http.ResponseWriter,
		req *http.Request,
	)
}

// NewHandler returns an initialized Handler instance
func NewHandler(
	core core.Core,
) Handler {
	return _handler{
		handleGetOrHeader: newHandleGetOrHeader(core),
		core:              core,
	}
}

type _handler struct {
	handleGetOrHeader handleGetOrHeader
	core              core.Core
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
	case "":
		http.Error(httpResp, "", http.StatusNotFound)
		return
	default:
		var pullCreds *model.PullCreds
		pullUsername, pullPassword, hasBasicAuth := httpReq.BasicAuth()
		if hasBasicAuth {
			pullCreds = &model.PullCreds{
				Username: pullUsername,
				Password: pullPassword,
			}
		}

		dataHandle, err := hdlr.core.ResolveData(
			httpReq.Context(),
			pathSegment,
			pullCreds,
		)
		if nil != err {
			var status int
			switch err.(type) {
			case model.ErrDataProviderAuthentication:
				hdlr.setWWWAuthenticateHeader(pathSegment, httpResp.Header())
				status = http.StatusUnauthorized
			case model.ErrDataProviderAuthorization:
				hdlr.setWWWAuthenticateHeader(pathSegment, httpResp.Header())
				status = http.StatusForbidden
			case model.ErrDataRefResolution:
				status = http.StatusNotFound
			default:
				status = http.StatusInternalServerError
			}
			http.Error(httpResp, err.Error(), status)
			return
		}

		hdlr.handleGetOrHeader.HandleGetOrHead(dataHandle, httpResp, httpReq)
	}

}

func (hdlr _handler) setWWWAuthenticateHeader(
	dataRef string,
	header http.Header,
) {
	realm := strings.SplitN(dataRef, "/", 2)[0]
	header.Set("WWW-Authenticate", fmt.Sprintf(`Basic realm="%s"`, realm))
}

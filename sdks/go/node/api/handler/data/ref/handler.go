package ref

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/opctl/opctl/sdks/go/internal/urlpath"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/core"
)

//counterfeiter:generate -o fakes/handler.go . Handler
type Handler interface {
	Handle(
		dataRef string,
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
	dataRef string,
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
		var pullCreds *model.Creds
		pullUsername, pullPassword, hasBasicAuth := httpReq.BasicAuth()
		if hasBasicAuth {
			pullCreds = &model.Creds{
				Username: pullUsername,
				Password: pullPassword,
			}
		}

		dataHandle, err := hdlr.core.ResolveData(
			httpReq.Context(),
			dataRef,
			pullCreds,
		)
		if err != nil {
			var status int
			if errors.Is(err, model.ErrDataProviderAuthentication{}) {
				hdlr.setWWWAuthenticateHeader(dataRef, httpResp.Header())
				status = http.StatusUnauthorized
			} else if errors.Is(err, model.ErrDataProviderAuthorization{}) {
				hdlr.setWWWAuthenticateHeader(dataRef, httpResp.Header())
				status = http.StatusForbidden
			} else if errors.Is(err, model.ErrDataRefResolution{}) {
				status = http.StatusNotFound
			} else {
				status = http.StatusInternalServerError
			}
			http.Error(httpResp, err.Error(), status)
			return
		}

		hdlr.handleGetOrHeader.HandleGetOrHead(dataHandle, httpResp, httpReq)
	default:
		http.Error(httpResp, "", http.StatusNotFound)
		return
	}
}

func (hdlr _handler) setWWWAuthenticateHeader(
	dataRef string,
	header http.Header,
) {
	realm := strings.SplitN(dataRef, "/", 2)[0]
	header.Set("WWW-Authenticate", fmt.Sprintf(`Basic realm="%s"`, realm))
}

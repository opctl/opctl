package liveness

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -o ./fakeHandler.go --fake-name FakeHandler ./ Handler

import (
	"github.com/opctl/opctl/sdks/go/node/core"
	"github.com/opctl/opctl/sdks/go/util/urlpath"
	"net/http"
)

type Handler interface {
	Handle(
		httpResp http.ResponseWriter,
		httpReq *http.Request,
	)
}

// NewHandler returns an initialized Handler instance
func NewHandler(
	core core.Core,
) Handler {
	return _handler{}
}

type _handler struct {
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

	if "" != httpReq.URL.Path || pathSegment != "liveness" {
		http.NotFoundHandler().ServeHTTP(httpResp, httpReq)
		return
	}

	// currently this is a no-op
}

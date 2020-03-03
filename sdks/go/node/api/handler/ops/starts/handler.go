package starts

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

import (
	"encoding/json"
	"net/http"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/core"
)

//counterfeiter:generate -o fakes/handler.go . Handler
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
		core: core,
	}
}

type _handler struct {
	core core.Core
}

func (hdlr _handler) Handle(
	httpResp http.ResponseWriter,
	httpReq *http.Request,
) {

	startOpReq := model.StartOpReq{}

	err := json.NewDecoder(httpReq.Body).Decode(&startOpReq)
	if nil != err {
		http.Error(httpResp, err.Error(), http.StatusBadRequest)
		return
	}

	callID, err := hdlr.core.StartOp(
		httpReq.Context(),
		startOpReq,
	)
	if nil != err {
		http.Error(httpResp, err.Error(), http.StatusInternalServerError)
		return
	}

	httpResp.WriteHeader(http.StatusCreated)
	httpResp.Header().Set("Content-Type", "text/plain; charset=UTF-8")

	httpResp.Write([]byte(callID))
	if nil != err {
		http.Error(httpResp, err.Error(), http.StatusInternalServerError)
		return
	}

}

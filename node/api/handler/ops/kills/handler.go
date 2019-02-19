package kills

//go:generate counterfeiter -o ./fakeHandler.go --fake-name FakeHandler ./ Handler

import (
	"encoding/json"
	"github.com/opctl/sdk-golang/model"
	"github.com/opctl/sdk-golang/node/core"
	"net/http"
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
	killOpReq := model.KillOpReq{}

	err := json.NewDecoder(httpReq.Body).Decode(&killOpReq)
	if nil != err {
		http.Error(httpResp, err.Error(), http.StatusBadRequest)
		return
	}

	hdlr.core.KillOp(killOpReq)

	httpResp.WriteHeader(http.StatusCreated)
	httpResp.Header().Set("Content-Type", "text/plain; charset=UTF-8")

}

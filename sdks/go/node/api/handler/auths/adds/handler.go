package adds

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

import (
	"encoding/json"
	"net/http"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node"
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
	node node.Node,
) Handler {
	return _handler{
		node: node,
	}
}

type _handler struct {
	node node.Node
}

func (hdlr _handler) Handle(
	httpResp http.ResponseWriter,
	httpReq *http.Request,
) {
	addAuthReq := model.AddAuthReq{}

	err := json.NewDecoder(httpReq.Body).Decode(&addAuthReq)
	if err != nil {
		http.Error(httpResp, err.Error(), http.StatusBadRequest)
		return
	}

	hdlr.node.AddAuth(httpReq.Context(), addAuthReq)

	httpResp.WriteHeader(http.StatusCreated)
	httpResp.Header().Set("Content-Type", "text/plain; charset=UTF-8")

}

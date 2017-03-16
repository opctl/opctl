package tcp

import (
	"encoding/json"
	"github.com/opspec-io/opctl/node/core"
	"github.com/opspec-io/sdk-golang/model"
	"net/http"
)

func newStartOpHandler(
	core core.Core,
) http.Handler {

	return &startOpHandler{
		core: core,
	}

}

type startOpHandler struct {
	core core.Core
}

func (this startOpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	startOpReq := model.StartOpReq{}

	err := json.NewDecoder(r.Body).Decode(&startOpReq)
	if nil != err {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	callId, err := this.core.StartOp(startOpReq)
	if nil != err {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")

	w.Write([]byte(callId))
	if nil != err {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

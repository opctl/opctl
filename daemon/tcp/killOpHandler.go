package tcp

import (
	"encoding/json"
	"github.com/opspec-io/engine/daemon/core"
	"github.com/opspec-io/sdk-golang/pkg/model"
	"net/http"
)

func newKillOpHandler(
	core core.Core,
) http.Handler {

	return &killOpHandler{
		core: core,
	}

}

type killOpHandler struct {
	core core.Core
}

func (this killOpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	killOpReq := model.KillOpReq{}

	err := json.NewDecoder(r.Body).Decode(&killOpReq)
	if nil != err {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	this.core.KillOp(killOpReq)

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")

}

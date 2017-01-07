package tcp

import (
	"github.com/gorilla/mux"
	"github.com/opspec-io/engine/daemon/core"
	"net/http"
)

type Api interface {
	Start()
}

func New(
	coreApi core.Core,
) Api {

	return &_api{
		compositionRoot: newCompositionRoot(coreApi),
	}

}

type _api struct {
	compositionRoot compositionRoot
}

func (this _api) Start() {

	router := mux.NewRouter()

	router.Handle(
		getEventBusRelUrlTemplate,
		this.compositionRoot.GetEventBusHandler(),
	).Methods(http.MethodGet)

	router.Handle(
		getLivenessRelUrlTemplate,
		this.compositionRoot.GetLivenessHandler(),
	).Methods(http.MethodGet)

	router.Handle(
		killOpRelUrlTemplate,
		this.compositionRoot.KillOpHandler(),
	).Methods(http.MethodPost)

	router.Handle(
		startOpRelUrlTemplate,
		this.compositionRoot.StartOpHandler(),
	).Methods(http.MethodPost)

	http.ListenAndServe(":42224", router)

}

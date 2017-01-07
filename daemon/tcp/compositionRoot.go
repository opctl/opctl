package tcp

import (
	"github.com/opspec-io/engine/daemon/core"
	"net/http"
)

type compositionRoot interface {
	GetEventBusHandler() http.Handler
	GetLivenessHandler() http.Handler
	KillOpHandler() http.Handler
	StartOpHandler() http.Handler
}

func newCompositionRoot(
	coreApi core.Core,
) (compositionRoot compositionRoot) {

	compositionRoot = &_compositionRoot{
		getEventBusHandler: newGetEventBusHandler(coreApi),
		getLivenessHandler: newGetLivenessHandler(),
		killOpHandler:      newKillOpHandler(coreApi),
		startOpHandler:     newStartOpHandler(coreApi),
	}

	return

}

type _compositionRoot struct {
	getLivenessHandler http.Handler
	getEventBusHandler http.Handler
	killOpHandler      http.Handler
	startOpHandler     http.Handler
}

func (this _compositionRoot) GetLivenessHandler() http.Handler {
	return this.getLivenessHandler
}

func (this _compositionRoot) GetEventBusHandler() http.Handler {
	return this.getEventBusHandler
}

func (this _compositionRoot) KillOpHandler() http.Handler {
	return this.killOpHandler
}

func (this _compositionRoot) StartOpHandler() http.Handler {
	return this.startOpHandler
}

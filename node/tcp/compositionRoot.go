package tcp

import (
	"github.com/opctl/opctl/node/core"
	"net/http"
)

type compositionRoot interface {
	GetEventStreamHandler() http.Handler
	GetLivenessHandler() http.Handler
	KillOpHandler() http.Handler
	StartOpHandler() http.Handler
}

func newCompositionRoot(
	coreApi core.Core,
) (compositionRoot compositionRoot) {

	compositionRoot = &_compositionRoot{
		getEventStreamHandler: newGetEventStreamHandler(coreApi),
		getLivenessHandler:    newGetLivenessHandler(),
		killOpHandler:         newKillOpHandler(coreApi),
		startOpHandler:        newStartOpHandler(coreApi),
	}

	return

}

type _compositionRoot struct {
	getLivenessHandler    http.Handler
	getEventStreamHandler http.Handler
	killOpHandler         http.Handler
	startOpHandler        http.Handler
}

func (this _compositionRoot) GetLivenessHandler() http.Handler {
	return this.getLivenessHandler
}

func (this _compositionRoot) GetEventStreamHandler() http.Handler {
	return this.getEventStreamHandler
}

func (this _compositionRoot) KillOpHandler() http.Handler {
	return this.killOpHandler
}

func (this _compositionRoot) StartOpHandler() http.Handler {
	return this.startOpHandler
}

package core

import (
	"github.com/opspec-io/opctl/util/uniquestring"
	"github.com/opspec-io/sdk-golang/pkg/model"
)

type serialOrchestrator interface {
	Execute(
		args map[string]*model.Arg,
		opGraphId string,
		opRef string,
		serialCall []*model.CallGraph,
	) (
		err error,
	)
}

func newSerialOrchestrator(
	orchestrator orchestrator,
	uniqueStringFactory uniquestring.UniqueStringFactory,
) serialOrchestrator {

	return &_serialOrchestrator{
		orchestrator:        orchestrator,
		uniqueStringFactory: uniqueStringFactory,
	}

}

type _serialOrchestrator struct {
	orchestrator        orchestrator
	uniqueStringFactory uniquestring.UniqueStringFactory
}

func (this _serialOrchestrator) Execute(
	args map[string]*model.Arg,
	opGraphId string,
	opRef string,
	serialCall []*model.CallGraph,
) (
	err error,
) {

	for _, childCall := range serialCall {
		err = this.orchestrator.Execute(
			this.uniqueStringFactory.Construct(),
			args,
			childCall,
			opRef,
			opGraphId,
		)
		if nil != err {
			// end run immediately on any error
			return
		}

	}

	return

}

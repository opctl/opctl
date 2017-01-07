package core

import (
	"github.com/opspec-io/opctl/util/uniquestring"
	"github.com/opspec-io/sdk-golang/pkg/model"
	"github.com/pkg/errors"
	"sync"
)

type parallelOrchestrator interface {
	Execute(
		args map[string]*model.Arg,
		opGraphId string,
		opRef string,
		parallelCall []*model.CallGraph,
	) (
		err error,
	)
}

func newParallelOrchestrator(
	orchestrator orchestrator,
	uniqueStringFactory uniquestring.UniqueStringFactory,
) parallelOrchestrator {

	return &_parallelOrchestrator{
		orchestrator:        orchestrator,
		uniqueStringFactory: uniqueStringFactory,
	}

}

type _parallelOrchestrator struct {
	orchestrator        orchestrator
	uniqueStringFactory uniquestring.UniqueStringFactory
}

func (this _parallelOrchestrator) Execute(
	args map[string]*model.Arg,
	opGraphId string,
	opRef string,
	parallelCall []*model.CallGraph,
) (
	err error,
) {

	var wg sync.WaitGroup
	var isSubOpRunErrors bool

	// run sub ops
	for _, childCall := range parallelCall {
		wg.Add(1)

		var childCallError error

		go func(childCall *model.CallGraph) {
			err = this.orchestrator.Execute(
				this.uniqueStringFactory.Construct(),
				args,
				childCall,
				opRef,
				opGraphId,
			)
			if nil != childCallError {
				isSubOpRunErrors = true
			}

			defer wg.Done()
		}(childCall)
	}
	wg.Wait()

	if isSubOpRunErrors {
		err = errors.New("One or more errors encountered in parallel run block")
	}

	return

}

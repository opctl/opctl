package core

//go:generate counterfeiter -o ./fakeParallelCaller.go --fake-name fakeParallelCaller ./ parallelCaller

import (
	"errors"
	"github.com/opspec-io/opctl/util/uniquestring"
	"github.com/opspec-io/sdk-golang/pkg/model"
	"sync"
)

type parallelCaller interface {
	// Executes a parallel call
	Call(
		inboundScope map[string]*model.Data,
		opGraphId string,
		opRef string,
		scgParallelCall []*model.Scg,
	) (
		err error,
	)
}

func newParallelCaller(
	caller caller,
	uniqueStringFactory uniquestring.UniqueStringFactory,
) parallelCaller {

	return _parallelCaller{
		caller:              caller,
		uniqueStringFactory: uniqueStringFactory,
	}

}

type _parallelCaller struct {
	caller              caller
	uniqueStringFactory uniquestring.UniqueStringFactory
}

func (this _parallelCaller) Call(
	inboundScope map[string]*model.Data,
	opGraphId string,
	opRef string,
	scgParallelCall []*model.Scg,
) (
	err error,
) {

	var wg sync.WaitGroup
	childErrChannel := make(chan error, len(scgParallelCall))

	for _, childCall := range scgParallelCall {
		wg.Add(1)

		go func(childCall *model.Scg) {
			// @TODO: handle sockets
			_, childErr := this.caller.Call(
				this.uniqueStringFactory.Construct(),
				inboundScope,
				childCall,
				opRef,
				opGraphId,
			)
			if nil != childErr {
				childErrChannel <- childErr
			}
			defer wg.Done()
		}(childCall)
	}
	wg.Wait()

	if len(childErrChannel) > 0 {
		// @TODO: consider including actual errors
		err = errors.New("One or more errors encountered in parallel run block")
	}

	return

}

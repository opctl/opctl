package core

//go:generate counterfeiter -o ./fakeParallelCaller.go --fake-name fakeParallelCaller ./ parallelCaller

import (
	"errors"
	"github.com/opspec-io/opctl/util/uniquestring"
	"github.com/opspec-io/sdk-golang/pkg/model"
	"sync"
)

type parallelCaller interface {
	Call(
		parentScope map[string]*model.Data,
		opGraphId string,
		opRef string,
		parallelCall []*model.Scg,
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
	parentScope map[string]*model.Data,
	opGraphId string,
	opRef string,
	parallelCall []*model.Scg,
) (
	err error,
) {

	var wg sync.WaitGroup
	childCallErrorChannel := make(chan error, len(parallelCall))

	for _, childCall := range parallelCall {
		wg.Add(1)

		go func(childCall *model.Scg) {
			wg.Done()
			// @TODO: handle sockets
			_, childCallErr := this.caller.Call(
				this.uniqueStringFactory.Construct(),
				parentScope,
				childCall,
				opRef,
				opGraphId,
			)
			if nil != childCallErr {
				childCallErrorChannel <- childCallErr
			}
		}(childCall)
	}
	wg.Wait()

	if len(childCallErrorChannel) > 0 {
		// @TODO: consider including actual errors
		err = errors.New("One or more errors encountered in parallel run block")
	}

	return

}

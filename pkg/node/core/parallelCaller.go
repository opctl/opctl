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
		inputs chan *variable,
		rootOpId string,
		pkgRef string,
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
	inputs chan *variable,
	rootOpId string,
	pkgRef string,
	scgParallelCall []*model.Scg,
) (
	err error,
) {

	// @TODO: stream in realtime
	inputMap := map[string]*model.Data{}
	for input := range inputs {
		inputMap[input.Name] = input.Value
	}

	var wg sync.WaitGroup
	childErrChannel := make(chan error, len(scgParallelCall))

	for _, childCall := range scgParallelCall {
		inputs := make(chan *variable, 150)
		for varName, varValue := range inputMap {
			inputs <- &variable{Name: varName, Value: varValue}
		}
		close(inputs)
		wg.Add(1)

		go func(childCall *model.Scg) {
			childErr := this.caller.Call(
				this.uniqueStringFactory.Construct(),
				// @TODO: broadcast output chan's to input chan's
				inputs,
				make(chan *variable, 150),
				childCall,
				pkgRef,
				rootOpId,
			)
			if nil != childErr {
				childErrChannel <- childErr
			}
			defer wg.Done()
		}(childCall)
	}
	wg.Wait()

	if len(childErrChannel) > 0 {
		err = errors.New("One or more errors encountered in parallel run block")
	}

	return

}

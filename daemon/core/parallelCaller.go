package core

//go:generate counterfeiter -o ./fakeParallelCaller.go --fake-name fakeParallelCaller ./ parallelCaller

import (
	"github.com/opspec-io/opctl/util/uniquestring"
	"github.com/opspec-io/sdk-golang/pkg/model"
	"github.com/pkg/errors"
	"sync"
)

type parallelCaller interface {
	Call(
		args map[string]*model.Data,
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

	return &_parallelCaller{
		caller:              caller,
		uniqueStringFactory: uniqueStringFactory,
	}

}

type _parallelCaller struct {
	caller              caller
	uniqueStringFactory uniquestring.UniqueStringFactory
}

func (this _parallelCaller) Call(
	args map[string]*model.Data,
	opGraphId string,
	opRef string,
	parallelCall []*model.Scg,
) (
	err error,
) {

	var wg sync.WaitGroup
	var isSubOpRunErrors bool

	// run sub ops
	for _, childCall := range parallelCall {
		wg.Add(1)

		var childCallError error

		go func(childCall *model.Scg) {
			// @TODO: handle sockets
			_, err = this.caller.Call(
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

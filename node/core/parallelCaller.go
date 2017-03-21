package core

//go:generate counterfeiter -o ./fakeParallelCaller.go --fake-name fakeParallelCaller ./ parallelCaller

import (
	"errors"
	"github.com/opctl/opctl/util/pubsub"
	"github.com/opctl/opctl/util/uniquestring"
	"github.com/opspec-io/sdk-golang/model"
	"sync"
	"time"
)

type parallelCaller interface {
	// Executes a parallel call
	Call(
		callId string,
		inboundScope map[string]*model.Data,
		rootOpId string,
		pkgRef string,
		scgParallelCall []*model.Scg,
	) (
		err error,
	)
}

func newParallelCaller(
	caller caller,
	pubSub pubsub.PubSub,
	uniqueStringFactory uniquestring.UniqueStringFactory,
) parallelCaller {

	return _parallelCaller{
		caller:              caller,
		pubSub:              pubSub,
		uniqueStringFactory: uniqueStringFactory,
	}

}

type _parallelCaller struct {
	caller              caller
	pubSub              pubsub.PubSub
	uniqueStringFactory uniquestring.UniqueStringFactory
}

func (this _parallelCaller) Call(
	callId string,
	inboundScope map[string]*model.Data,
	rootOpId string,
	pkgRef string,
	scgParallelCall []*model.Scg,
) (
	err error,
) {

	defer func() {
		// defer must be defined before conditional return statements so it always runs

		this.pubSub.Publish(
			&model.Event{
				Timestamp: time.Now().UTC(),
				ParallelCallEnded: &model.ParallelCallEndedEvent{
					CallId:   callId,
					RootOpId: rootOpId,
				},
			},
		)

	}()

	var wg sync.WaitGroup
	childErrChannel := make(chan error, len(scgParallelCall))

	for _, childCall := range scgParallelCall {
		wg.Add(1)

		go func(childCall *model.Scg) {
			childErr := this.caller.Call(
				this.uniqueStringFactory.Construct(),
				inboundScope,
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

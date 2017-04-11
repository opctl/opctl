package core

//go:generate counterfeiter -o ./fakeParallelCaller.go --fake-name fakeParallelCaller ./ parallelCaller

import (
	"bytes"
	"fmt"
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
		scgParallelCall []*model.SCG,
	) (
		err error,
	)
}

func newParallelCaller(
	caller caller,
	opKiller opKiller,
	pubSub pubsub.PubSub,
	uniqueStringFactory uniquestring.UniqueStringFactory,
) parallelCaller {

	return _parallelCaller{
		opKiller:            opKiller,
		caller:              caller,
		pubSub:              pubSub,
		uniqueStringFactory: uniqueStringFactory,
	}

}

type _parallelCaller struct {
	opKiller            opKiller
	caller              caller
	pubSub              pubsub.PubSub
	uniqueStringFactory uniquestring.UniqueStringFactory
}

func (this _parallelCaller) Call(
	callId string,
	inboundScope map[string]*model.Data,
	rootOpId string,
	pkgRef string,
	scgParallelCall []*model.SCG,
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
	childErrChannel := make(chan error, 1)

	// setup cancellation
	cancellationChannel := make(chan struct{})
	cancellationRequestChannel := make(chan error, len(scgParallelCall))
	go func() {
		// process cancellation requests
		childErr := <-cancellationRequestChannel

		// record error
		childErrChannel <- childErr

		close(cancellationChannel)
	}()

	// perform calls in parallel w/ cancellation
	for _, childCall := range scgParallelCall {
		wg.Add(1)

		go func(childCall *model.SCG) {
			defer wg.Done()

			childDoneChannel := make(chan struct{})
			go func() {
				defer close(childDoneChannel)
				childErr := this.caller.Call(
					this.uniqueStringFactory.Construct(),
					inboundScope,
					childCall,
					pkgRef,
					rootOpId,
				)
				if nil != childErr {
					cancellationRequestChannel <- childErr
				}
			}()

			select {
			case <-cancellationChannel:
				// ensure resources immediately reclaimed
				this.opKiller.Kill(model.KillOpReq{OpId: rootOpId})
			case <-childDoneChannel:
			}
		}(childCall)
	}
	wg.Wait()

	if len(childErrChannel) == 0 {
		// don't leak go routine
		close(cancellationRequestChannel)
	} else {

		messageBuffer := bytes.NewBufferString(
			fmt.Sprint(`
-
  Error during parallel call.
  Error:`))
		childErr := <-childErrChannel
		messageBuffer.WriteString(fmt.Sprintf(`
    - %v`,
			childErr.Error(),
		))
		err = fmt.Errorf(
			`%v
-`, messageBuffer.String())
	}

	return

}

package core

//go:generate counterfeiter -o ./fakeParallelCaller.go --fake-name fakeParallelCaller ./ parallelCaller

import (
	"bytes"
	"fmt"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/util/pubsub"
	"github.com/opspec-io/sdk-golang/util/uniquestring"
	"sync"
	"time"
)

type parallelCaller interface {
	// Executes a parallel call
	Call(
		callId string,
		inboundScope map[string]*model.Value,
		rootOpId string,
		pkgRef string,
		scgParallelCall []*model.SCG,
	) error
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
	inboundScope map[string]*model.Value,
	rootOpId string,
	pkgRef string,
	scgParallelCall []*model.SCG,
) error {

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

	// setup cancellation
	cancellationChannel := make(chan struct{})
	cancellationReqChannel := make(chan struct{}, len(scgParallelCall))
	go func() {
		<-cancellationReqChannel
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
				if childErr := this.caller.Call(
					this.uniqueStringFactory.Construct(),
					inboundScope,
					childCall,
					pkgRef,
					rootOpId,
				); nil != childErr {
					childErrChannel <- childErr
					// trigger cancellation
					cancellationReqChannel <- struct{}{}
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
		close(cancellationReqChannel)
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
		return fmt.Errorf(
			`%v
-`, messageBuffer.String())
	}

	return nil

}

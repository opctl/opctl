package core

//go:generate counterfeiter -o ./fakeParallelCaller.go --fake-name fakeParallelCaller ./ parallelCaller

import (
	"bytes"
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/opctl/sdk-golang/model"
	"github.com/opctl/sdk-golang/util/pubsub"
	"github.com/opctl/sdk-golang/util/uniquestring"
)

type parallelCaller interface {
	// Executes a parallel call
	Call(
		ctx context.Context,
		callId string,
		inboundScope map[string]*model.Value,
		rootOpID string,
		opHandle model.DataHandle,
		scgParallelCall []*model.SCG,
	) error
}

func newParallelCaller(
	caller caller,
	callKiller callKiller,
	pubSub pubsub.PubSub,
) parallelCaller {

	return _parallelCaller{
		callKiller:          callKiller,
		caller:              caller,
		pubSub:              pubSub,
		uniqueStringFactory: uniquestring.NewUniqueStringFactory(),
	}

}

type _parallelCaller struct {
	callKiller          callKiller
	caller              caller
	pubSub              pubsub.PubSub
	uniqueStringFactory uniquestring.UniqueStringFactory
}

func (this _parallelCaller) Call(
	ctx context.Context,
	callId string,
	inboundScope map[string]*model.Value,
	rootOpID string,
	opHandle model.DataHandle,
	scgParallelCall []*model.SCG,
) error {

	defer func() {

		// defer must be defined before conditional return statements so it always runs
		this.pubSub.Publish(
			model.Event{
				Timestamp: time.Now().UTC(),
				ParallelCallEnded: &model.ParallelCallEndedEvent{
					CallID:   callId,
					RootOpID: rootOpID,
				},
			},
		)

	}()

	var wg sync.WaitGroup
	childErrChannel := make(chan error, len(scgParallelCall))

	// setup cancellation
	parallelCtx, parallelCancel := context.WithCancel(ctx)
	defer parallelCancel()

	// perform calls in parallel w/ cancellation
	for _, childCall := range scgParallelCall {
		wg.Add(1)

		go func(childCall *model.SCG) {
			defer wg.Done()

			childCallID, err := this.uniqueStringFactory.Construct()
			if nil != err {
				childErrChannel <- err
				// trigger parallel cancellation
				parallelCancel()
			}

			childCtx, childCancel := context.WithCancel(parallelCtx)
			go func() {
				defer childCancel()
				if childErr := this.caller.Call(
					childCtx,
					childCallID,
					inboundScope,
					childCall,
					opHandle,
					&callId,
					rootOpID,
				); nil != childErr {
					childErrChannel <- childErr
					// trigger parallel cancellation
					parallelCancel()
				}
			}()

			select {
			case <-parallelCtx.Done():
				// ensure resources immediately reclaimed
				this.callKiller.Kill(
					childCallID,
					rootOpID,
				)
			case <-childCtx.Done():
			}
		}(childCall)
	}
	wg.Wait()

	if len(childErrChannel) != 0 {
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

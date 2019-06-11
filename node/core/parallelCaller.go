package core

//go:generate counterfeiter -o ./fakeParallelCaller.go --fake-name fakeParallelCaller ./ parallelCaller

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/opctl/sdk-golang/model"
	"github.com/opctl/sdk-golang/util/pubsub"
	"github.com/opctl/sdk-golang/util/uniquestring"
)

type parallelCaller interface {
	// Executes a parallel call
	Call(
		ctx context.Context,
		callID string,
		inboundScope map[string]*model.Value,
		rootOpID string,
		opHandle model.DataHandle,
		scgParallelCall []*model.SCG,
	)
}

func newParallelCaller(
	caller caller,
	pubSub pubsub.PubSub,
) parallelCaller {

	return _parallelCaller{
		caller:              caller,
		pubSub:              pubSub,
		uniqueStringFactory: uniquestring.NewUniqueStringFactory(),
	}

}

type _parallelCaller struct {
	caller              caller
	pubSub              pubsub.PubSub
	uniqueStringFactory uniquestring.UniqueStringFactory
}

func (pc _parallelCaller) Call(
	ctx context.Context,
	callID string,
	inboundScope map[string]*model.Value,
	rootOpID string,
	opHandle model.DataHandle,
	scgParallelCall []*model.SCG,
) {
	// setup cancellation
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	outputs := map[string]*model.Value{}
	var err error

	defer func() {
		// defer must be defined before conditional return statements so it always runs
		event := model.Event{
			ParallelCallEnded: &model.ParallelCallEndedEvent{
				CallID:   callID,
				Outputs:  outputs,
				RootOpID: rootOpID,
			},
			Timestamp: time.Now().UTC(),
		}

		if nil != err {
			event.ParallelCallEnded.Error = &model.CallEndedEventError{
				Message: err.Error(),
			}
		}
		pc.pubSub.Publish(
			event,
		)
	}()

	childCallIDIndexMap := map[string]int{}
	callIndexOutputsMap := map[int]map[string]*model.Value{}

	// perform calls in parallel w/ cancellation
	for childCallIndex, childCall := range scgParallelCall {

		var childCallID string
		childCallID, err = pc.uniqueStringFactory.Construct()
		if nil != err {
			// trigger parallel cancellation
			cancel()
		}
		childCallIDIndexMap[childCallID] = childCallIndex

		go pc.caller.Call(
			ctx,
			childCallID,
			inboundScope,
			childCall,
			opHandle,
			&callID,
			rootOpID,
		)
	}

	// subscribe to events
	// @TODO: handle err channel
	eventFilterSince := time.Now().UTC()
	eventChannel, _ := pc.pubSub.Subscribe(
		ctx,
		model.EventFilter{
			Roots: []string{rootOpID},
			Since: &eventFilterSince,
		},
	)

	childErrorMessages := []string{}
eventLoop:
	for event := range eventChannel {
		if nil != event.CallEnded {
			if childCallIndex, isChildCallEnded := childCallIDIndexMap[event.CallEnded.CallID]; isChildCallEnded {
				callIndexOutputsMap[childCallIndex] = event.CallEnded.Outputs
				if nil != event.CallEnded.Error {
					cancel()
					childErrorMessages = append(childErrorMessages, event.CallEnded.Error.Message)
				}
			}

			if len(callIndexOutputsMap) == len(scgParallelCall) {
				// construct outputs
				for i := 0; i < len(scgParallelCall); i++ {
					callOutputs := callIndexOutputsMap[i]
					for varName, varData := range callOutputs {
						outputs[varName] = varData
					}
				}

				break eventLoop
			}

		}
	}

	if len(childErrorMessages) != 0 {
		err = fmt.Errorf(
			"-\nError(s) during parallel call. Error(s) were:\n%v\n-",
			strings.Join(childErrorMessages, "\n"),
		)
	}

}

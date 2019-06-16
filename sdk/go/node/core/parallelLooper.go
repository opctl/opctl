package core

//go:generate counterfeiter -o ./fakeLooper.go --fake-name fakeLooper ./ looper

import (
	"context"
	"fmt"
	"time"

	"github.com/opctl/opctl/sdk/go/opspec/interpreter/call/loop"
	"github.com/opctl/opctl/sdk/go/opspec/interpreter/call/loop/iteration"
	"github.com/opctl/opctl/sdk/go/opspec/interpreter/loopable"

	"github.com/opctl/opctl/sdk/go/model"
	"github.com/opctl/opctl/sdk/go/util/pubsub"
	"github.com/opctl/opctl/sdk/go/util/uniquestring"
)

type parallelLooper interface {
	// Loop loops a call
	Loop(
		ctx context.Context,
		id string,
		inboundScope map[string]*model.Value,
		scg *model.SCG,
		opHandle model.DataHandle,
		parentCallID *string,
		rootOpID string,
	)
}

func newParallelLooper(
	caller caller,
	pubSub pubsub.PubSub,
) looper {
	return _parallelLooper{
		caller:              caller,
		iterationScoper:     iteration.NewScoper(),
		loopableInterpreter: loopable.NewInterpreter(),
		loopDeScoper:        loop.NewDeScoper(),
		pubSub:              pubSub,
		uniqueStringFactory: uniquestring.NewUniqueStringFactory(),
		loopInterpreter:     loop.NewInterpreter(),
	}
}

type _parallelLooper struct {
	caller              caller
	iterationScoper     iteration.Scoper
	loopableInterpreter loopable.Interpreter
	loopDeScoper        loop.DeScoper
	pubSub              pubsub.PubSub
	uniqueStringFactory uniquestring.UniqueStringFactory
	loopInterpreter     loop.Interpreter
}

func (plpr _parallelLooper) Loop(
	ctx context.Context,
	id string,
	inboundScope map[string]*model.Value,
	scg *model.SCG,
	opHandle model.DataHandle,
	parentCallID *string,
	rootOpID string,
) {
	// setup cancellation
	ctxOfChildren, cancelChildren := context.WithCancel(ctx)
	defer cancelChildren()

	outboundScope := map[string]*model.Value{}
	var err error

	defer func() {
		// defer must be defined before conditional return statements so it always runs
		event := model.Event{
			Timestamp: time.Now().UTC(),
			CallEnded: &model.CallEndedEvent{
				CallID:     id,
				RootCallID: rootOpID,
				Outputs:    outboundScope,
			},
		}

		if nil != err {
			event.CallEnded.Error = &model.CallEndedEventError{
				Message: err.Error(),
			}
		}

		plpr.pubSub.Publish(event)
	}()

	childCallIndex := 0
	outboundScope, err = plpr.iterationScoper.Scope(
		childCallIndex,
		inboundScope,
		scg.Loop,
		opHandle,
	)
	if nil != err {
		return
	}

	// copy scg.Loop & remove from scg since we're already looping
	scgLoop := scg.Loop
	scg.Loop = nil

	// interpret initial iteration of the loop
	var dcgLoop *model.DCGLoop
	dcgLoop, err = plpr.loopInterpreter.Interpret(
		opHandle,
		scgLoop,
		outboundScope,
	)
	if nil != err {
		return
	}

	startTime := time.Now().UTC()
	childCallIDIndexMap := map[string]int{}
	callIndexOutputsMap := map[int]map[string]*model.Value{}

	for !loop.IsIterationComplete(childCallIndex, dcgLoop) {

		var childCallID string
		childCallID, err = plpr.uniqueStringFactory.Construct()
		if nil != err {
			// cancel all children on any error
			cancelChildren()
		}
		childCallIDIndexMap[childCallID] = childCallIndex

		go plpr.caller.Call(
			ctxOfChildren,
			childCallID,
			outboundScope,
			scg,
			opHandle,
			parentCallID,
			rootOpID,
		)

		childCallIndex++

		if loop.IsIterationComplete(childCallIndex, dcgLoop) {
			break
		}

		outboundScope, err = plpr.iterationScoper.Scope(
			childCallIndex,
			outboundScope,
			scgLoop,
			opHandle,
		)
		if nil != err {
			return
		}

		// interpret next iteration of the loop
		dcgLoop, err = plpr.loopInterpreter.Interpret(
			opHandle,
			scgLoop,
			outboundScope,
		)
		if nil != err {
			return
		}

	}

	// subscribe to events
	// @TODO: handle err channel
	eventChannel, _ := plpr.pubSub.Subscribe(
		ctx,
		model.EventFilter{
			Roots: []string{rootOpID},
			Since: &startTime,
		},
	)

	if len(childCallIDIndexMap) == 0 {
		return
	}

	childErrorMessages := []string{}
	for event := range eventChannel {
		if nil != event.CallEnded {
			if childCallIndex, isChildCallEnded := childCallIDIndexMap[event.CallEnded.CallID]; isChildCallEnded {
				callIndexOutputsMap[childCallIndex] = event.CallEnded.Outputs
				if nil != event.CallEnded.Error {
					// cancel all children on any error
					cancelChildren()
					childErrorMessages = append(childErrorMessages, event.CallEnded.Error.Message)
				}
			}

			if len(callIndexOutputsMap) == len(childCallIDIndexMap) {
				// all calls have ended

				// construct parallel outputs
				for i := 0; i < len(childCallIDIndexMap); i++ {
					callOutputs := callIndexOutputsMap[i]
					for varName, varData := range callOutputs {
						outboundScope[varName] = varData
					}
				}

				// construct parallel error
				if len(childErrorMessages) != 0 {
					var formattedChildErrorMessages string
					for _, childErrorMessage := range childErrorMessages {
						formattedChildErrorMessages = fmt.Sprintf("\t-%v\n", childErrorMessage)
					}
					err = fmt.Errorf(
						"-\nError(s) during parallel call. Error(s) were:\n%v\n-",
						formattedChildErrorMessages,
					)
				}

				return
			}

		}
	}

	outboundScope = plpr.loopDeScoper.DeScope(
		inboundScope,
		scgLoop,
		outboundScope,
	)
}

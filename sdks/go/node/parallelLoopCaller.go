package node

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/loop"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/loop/iteration"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/parallelloop"

	"github.com/opctl/opctl/sdks/go/internal/uniquestring"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/pubsub"
)

//counterfeiter:generate -o internal/fakes/parallelLoopCaller.go . parallelLoopCaller
type parallelLoopCaller interface {
	// Executes a parallel loop call
	Call(
		parentCtx context.Context,
		id string,
		inboundScope map[string]*model.Value,
		callSpecParallelLoop model.ParallelLoopCallSpec,
		opPath string,
		parentCallID *string,
		rootCallID string,
	) (
		map[string]*model.Value,
		error,
	)
}

func newParallelLoopCaller(
	caller caller,
	pubSub pubsub.PubSub,
) parallelLoopCaller {
	return _parallelLoopCaller{
		caller: caller,
		pubSub: pubSub,
	}
}

type _parallelLoopCaller struct {
	caller caller
	pubSub pubsub.PubSub
}

func (plpr _parallelLoopCaller) Call(
	parentCtx context.Context,
	id string,
	inboundScope map[string]*model.Value,
	callSpecParallelLoop model.ParallelLoopCallSpec,
	opPath string,
	parentCallID *string,
	rootCallID string,
) (
	map[string]*model.Value,
	error,
) {
	// setup cancellation
	parallelLoopCtx, cancelParallelLoop := context.WithCancel(parentCtx)
	defer cancelParallelLoop()

	childCallIndex := 0
	startTime := time.Now().UTC()
	childCallIndexByID := map[string]int{}

	for {

		childCallID, err := uniquestring.Construct()
		if err != nil {
			// end run immediately on any error
			return nil, err
		}

		childCallScope, scopeErr := iteration.Scope(
			childCallIndex,
			inboundScope,
			callSpecParallelLoop.Range,
			callSpecParallelLoop.Vars,
		)
		if scopeErr != nil {
			return nil, scopeErr
		}

		// interpret iteration of the loop
		callParallelLoop, interpretErr := parallelloop.Interpret(
			callSpecParallelLoop,
			childCallScope,
		)
		if interpretErr != nil {
			return nil, interpretErr
		}

		if parallelloop.IsIterationComplete(childCallIndex, *callParallelLoop) {
			break
		}

		childCallIndexByID[childCallID] = childCallIndex

		go func() {
			defer func() {
				if panicArg := recover(); panicArg != nil {
					// recover from panics; treat as errors
					fmt.Printf("%v\n%v", panicArg, debug.Stack())

					// cancel all children on any error
					cancelParallelLoop()
				}
			}()

			plpr.caller.Call(
				parallelLoopCtx,
				childCallID,
				childCallScope,
				&callSpecParallelLoop.Run,
				opPath,
				parentCallID,
				rootCallID,
			)
		}()

		childCallIndex++

	}

	if len(childCallIndexByID) == 0 {
		return nil, nil
	}

	// subscribe to events
	// @TODO: handle err channel
	eventChannel, _ := plpr.pubSub.Subscribe(
		// don't cancel w/ children; we need to read err msgs
		parentCtx,
		model.EventFilter{
			Roots: []string{rootCallID},
			Since: &startTime,
		},
	)

	var isChildErred = false
	childCallOutputsByIndex := map[int]map[string]*model.Value{}
	outputs := inboundScope

eventLoop:
	for event := range eventChannel {
		if event.CallEnded != nil {
			if childCallIndex, isChildCallEnded := childCallIndexByID[event.CallEnded.Call.ID]; isChildCallEnded {
				childCallOutputsByIndex[childCallIndex] = event.CallEnded.Outputs
				if event.CallEnded.Error != nil {
					isChildErred = true

					// cancel all children on any error
					cancelParallelLoop()
				}
			}

			if len(childCallOutputsByIndex) == len(childCallIndexByID) {
				// all calls have ended

				// construct parallel outputs
				for i := 0; i < len(childCallIndexByID); i++ {
					callOutputs := childCallOutputsByIndex[i]
					for varName, varData := range callOutputs {
						outputs[varName] = varData
					}
				}

				if isChildErred {
					return nil, errors.New("child call failed")
				}

				break eventLoop
			}

		}
	}

	return loop.DeScope(
		inboundScope,
		callSpecParallelLoop.Range,
		callSpecParallelLoop.Vars,
		outputs,
	), nil
}

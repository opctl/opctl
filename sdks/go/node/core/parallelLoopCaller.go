package core

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/loop"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/loop/iteration"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/parallelloop"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/loopable"

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
		caller:                  caller,
		iterationScoper:         iteration.NewScoper(),
		loopableInterpreter:     loopable.NewInterpreter(),
		loopDeScoper:            loop.NewDeScoper(),
		parallelLoopInterpreter: parallelloop.NewInterpreter(),
		pubSub:                  pubSub,
		uniqueStringFactory:     uniquestring.NewUniqueStringFactory(),
	}
}

type _parallelLoopCaller struct {
	caller                  caller
	iterationScoper         iteration.Scoper
	loopableInterpreter     loopable.Interpreter
	loopDeScoper            loop.DeScoper
	parallelLoopInterpreter parallelloop.Interpreter
	pubSub                  pubsub.PubSub
	uniqueStringFactory     uniquestring.UniqueStringFactory
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
	outboundScope, scopeErr := plpr.iterationScoper.Scope(
		childCallIndex,
		inboundScope,
		callSpecParallelLoop.Range,
		callSpecParallelLoop.Vars,
	)
	if nil != scopeErr {
		return nil, scopeErr
	}

	// interpret initial iteration of the loop
	callParallelLoop, interpretErr := plpr.parallelLoopInterpreter.Interpret(
		callSpecParallelLoop,
		outboundScope,
	)
	if nil != interpretErr {
		return nil, interpretErr
	}

	startTime := time.Now().UTC()
	childCallIDIndexMap := map[string]int{}
	callIndexOutputsMap := map[int]map[string]*model.Value{}

	for !parallelloop.IsIterationComplete(childCallIndex, *callParallelLoop) {

		childCallID, err := plpr.uniqueStringFactory.Construct()
		if nil != err {
			// end run immediately on any error
			return nil, err
		}
		childCallIDIndexMap[childCallID] = childCallIndex

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
				outboundScope,
				&callSpecParallelLoop.Run,
				opPath,
				parentCallID,
				rootCallID,
			)
		}()

		childCallIndex++

		if parallelloop.IsIterationComplete(childCallIndex, *callParallelLoop) {
			break
		}

		var scopeErr error
		outboundScope, scopeErr = plpr.iterationScoper.Scope(
			childCallIndex,
			outboundScope,
			callSpecParallelLoop.Range,
			callSpecParallelLoop.Vars,
		)
		if nil != scopeErr {
			return nil, scopeErr
		}

		// interpret next iteration of the loop
		var interpretErr error
		callParallelLoop, interpretErr = plpr.parallelLoopInterpreter.Interpret(
			callSpecParallelLoop,
			outboundScope,
		)
		if nil != interpretErr {
			return nil, interpretErr
		}

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

	if len(childCallIDIndexMap) == 0 {
		return nil, nil
	}

	var isChildErred = false

eventLoop:
	for event := range eventChannel {
		if nil != event.CallEnded {
			if childCallIndex, isChildCallEnded := childCallIDIndexMap[event.CallEnded.Call.ID]; isChildCallEnded {
				callIndexOutputsMap[childCallIndex] = event.CallEnded.Outputs
				if nil != event.CallEnded.Error {
					isChildErred = true

					// cancel all children on any error
					cancelParallelLoop()
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

				if isChildErred {
					return nil, errors.New("child call failed")
				}

				break eventLoop
			}

		}
	}

	return plpr.loopDeScoper.DeScope(
		inboundScope,
		callSpecParallelLoop.Range,
		callSpecParallelLoop.Vars,
		outboundScope,
	), nil
}

package core

import (
	"context"
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
) {
	// setup cancellation
	parallelLoopCtx, cancelParallelLoop := context.WithCancel(parentCtx)
	defer cancelParallelLoop()

	outboundScope := map[string]*model.Value{}
	var err error

	defer func() {
		// defer must be defined before conditional return statements so it always runs
		event := model.Event{
			Timestamp: time.Now().UTC(),
			CallEnded: &model.CallEnded{
				CallID:     id,
				RootCallID: rootCallID,
				Outputs:    outboundScope,
			},
		}

		if nil != err {
			event.CallEnded.Error = &model.CallEndedError{
				Message: err.Error(),
			}
		}

		plpr.pubSub.Publish(event)
	}()

	childCallIndex := 0
	outboundScope, err = plpr.iterationScoper.Scope(
		childCallIndex,
		inboundScope,
		callSpecParallelLoop.Range,
		callSpecParallelLoop.Vars,
	)
	if nil != err {
		return
	}

	// interpret initial iteration of the loop
	var dcgParallelLoop *model.ParallelLoopCall
	dcgParallelLoop, err = plpr.parallelLoopInterpreter.Interpret(
		callSpecParallelLoop,
		outboundScope,
	)
	if nil != err {
		return
	}

	startTime := time.Now().UTC()
	childCallIDIndexMap := map[string]int{}
	callIndexOutputsMap := map[int]map[string]*model.Value{}

	for !parallelloop.IsIterationComplete(childCallIndex, *dcgParallelLoop) {

		var childCallID string
		childCallID, err = plpr.uniqueStringFactory.Construct()
		if nil != err {
			// cancel all children on any error
			cancelParallelLoop()
		}
		childCallIDIndexMap[childCallID] = childCallIndex

		go func() {
			defer func() {
				if panicArg := recover(); panicArg != nil {
					// recover from panics; treat as errors
					err = fmt.Errorf("%v\n%v", panicArg, debug.Stack())

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

		if parallelloop.IsIterationComplete(childCallIndex, *dcgParallelLoop) {
			break
		}

		outboundScope, err = plpr.iterationScoper.Scope(
			childCallIndex,
			outboundScope,
			callSpecParallelLoop.Range,
			callSpecParallelLoop.Vars,
		)
		if nil != err {
			return
		}

		// interpret next iteration of the loop
		dcgParallelLoop, err = plpr.parallelLoopInterpreter.Interpret(
			callSpecParallelLoop,
			outboundScope,
		)
		if nil != err {
			return
		}

	}

	// subscribe to events
	// @TODO: handle err channel
	eventChannel, _ := plpr.pubSub.Subscribe(
		parallelLoopCtx,
		model.EventFilter{
			Roots: []string{rootCallID},
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
					cancelParallelLoop()
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
		callSpecParallelLoop.Range,
		callSpecParallelLoop.Vars,
		outboundScope,
	)
}

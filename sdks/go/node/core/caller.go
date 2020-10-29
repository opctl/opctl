package core

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call"
	"github.com/opctl/opctl/sdks/go/pubsub"
)

//counterfeiter:generate -o internal/fakes/caller.go . caller
type caller interface {
	// Call executes a call
	Call(
		ctx context.Context,
		id string,
		scope map[string]*model.Value,
		callSpec *model.CallSpec,
		opPath string,
		parentCallID *string,
		rootCallID string,
	)
}

func newCaller(
	callInterpreter call.Interpreter,
	containerCaller containerCaller,
	dataDirPath string,
	callStore callStore,
	pubSub pubsub.PubSub,
) caller {
	instance := &_caller{
		callInterpreter: callInterpreter,
		containerCaller: containerCaller,
		callStore:       callStore,
		pubSub:          pubSub,
	}

	instance.opCaller = newOpCaller(
		callStore,
		pubSub,
		instance,
		dataDirPath,
	)

	instance.parallelCaller = newParallelCaller(
		instance,
		pubSub,
	)

	instance.parallelLoopCaller = newParallelLoopCaller(
		instance,
		pubSub,
	)

	instance.serialCaller = newSerialCaller(
		instance,
		pubSub,
	)

	instance.serialLoopCaller = newSerialLoopCaller(
		instance,
		pubSub,
	)

	return instance
}

type _caller struct {
	callInterpreter    call.Interpreter
	containerCaller    containerCaller
	callStore          callStore
	opCaller           opCaller
	parallelCaller     parallelCaller
	parallelLoopCaller parallelLoopCaller
	pubSub             pubsub.PubSub
	serialCaller       serialCaller
	serialLoopCaller   serialLoopCaller
}

func (clr _caller) Call(
	ctx context.Context,
	id string,
	scope map[string]*model.Value,
	callSpec *model.CallSpec,
	opPath string,
	parentCallID *string,
	rootCallID string,
) {
	callCtx, cancelCall := context.WithCancel(ctx)
	var err error
	var isKilled bool
	var outputs map[string]*model.Value
	callStartTime := time.Now().UTC()

	defer func() {
		// defer must be defined before conditional return statements so it always runs
		<-callCtx.Done()
		var outcome string
		if isKilled {
			outcome = model.OpOutcomeKilled
		} else if nil != err {
			outcome = model.OpOutcomeFailed
		} else {
			outcome = model.OpOutcomeSucceeded
		}

		event := model.Event{
			CallEnded: &model.CallEnded{
				CallID:     id,
				Outcome:    outcome,
				Outputs:    outputs,
				RootCallID: rootCallID,
			},
			Timestamp: time.Now().UTC(),
		}

		if outcome == model.OpOutcomeFailed {
			event.CallEnded.Error = &model.CallEndedError{
				Message: err.Error(),
			}
		}

		clr.pubSub.Publish(event)
	}()

	if nil == callSpec {
		cancelCall()

		// NOOP
		return
	}

	var dcg *model.Call
	dcg, err = clr.callInterpreter.Interpret(
		scope,
		callSpec,
		id,
		opPath,
		parentCallID,
		rootCallID,
	)
	if nil != err {
		cancelCall()

		return
	}

	if nil != dcg.If && !*dcg.If {
		cancelCall()

		return
	}

	go func() {
		defer func() {
			if panicArg := recover(); panicArg != nil {
				// recover from panics; treat as errors
				err = fmt.Errorf(
					fmt.Sprint(panicArg, debug.Stack()),
				)
			}
		}()

		defer cancelCall()

		eventChannel, _ := clr.pubSub.Subscribe(
			callCtx,
			model.EventFilter{
				Roots: []string{rootCallID},
				Since: &callStartTime,
			},
		)

		for event := range eventChannel {
			switch {
			case nil != event.CallEnded && event.CallEnded.CallID == id:
				if nil != event.CallEnded.Error {
					err = errors.New(event.CallEnded.Error.Message)
				}
				outputs = event.CallEnded.Outputs
				return
			case nil != event.ContainerExited && event.ContainerExited.ContainerID == id:
				if nil != event.ContainerExited.Error {
					err = errors.New(event.ContainerExited.Error.Message)
				}
				outputs = event.ContainerExited.Outputs
				return
			case nil != event.OpKillRequested:
				isKilled = true
				return
			}
		}
	}()

	clr.callStore.Add(dcg)

	switch {
	case nil != callSpec.Container:
		clr.containerCaller.Call(
			callCtx,
			dcg.Container,
			scope,
			callSpec.Container,
		)
	case nil != callSpec.Op:
		clr.opCaller.Call(
			callCtx,
			dcg.Op,
			scope,
			parentCallID,
			callSpec.Op,
		)
	case nil != callSpec.Parallel:
		clr.parallelCaller.Call(
			callCtx,
			id,
			scope,
			rootCallID,
			opPath,
			*callSpec.Parallel,
		)
	case nil != callSpec.ParallelLoop:
		clr.parallelLoopCaller.Call(
			callCtx,
			id,
			scope,
			*callSpec.ParallelLoop,
			opPath,
			parentCallID,
			rootCallID,
		)
	case nil != callSpec.Serial:
		clr.serialCaller.Call(
			callCtx,
			id,
			scope,
			rootCallID,
			opPath,
			*callSpec.Serial,
		)
	case nil != callSpec.SerialLoop:
		clr.serialLoopCaller.Call(
			callCtx,
			id,
			scope,
			*callSpec.SerialLoop,
			opPath,
			parentCallID,
			rootCallID,
		)
	default:
		err = fmt.Errorf("Invalid call graph %+v\n", callSpec)
	}

}

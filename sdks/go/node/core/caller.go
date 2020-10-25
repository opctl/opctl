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
		rootOpID string,
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
	rootOpID string,
) {
	callCtx, cancelCall := context.WithCancel(ctx)
	var err error
	var outputs map[string]*model.Value
	callStartTime := time.Now().UTC()

	defer func() {
		// defer must be defined before conditional return statements so it always runs
		<-callCtx.Done()

		event := model.Event{
			CallEnded: &model.CallEnded{
				CallID:     id,
				Outputs:    outputs,
				RootCallID: rootOpID,
			},
			Timestamp: time.Now().UTC(),
		}

		if nil != err {
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
		rootOpID,
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
				Roots: []string{rootOpID},
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
			case nil != event.OpEnded && event.OpEnded.OpID == id:
				if nil != event.OpEnded.Error {
					err = errors.New(event.OpEnded.Error.Message)
				}
				outputs = event.OpEnded.Outputs
				return
			case nil != event.SerialCallEnded && event.SerialCallEnded.CallID == id:
				if nil != event.SerialCallEnded.Error {
					err = errors.New(event.SerialCallEnded.Error.Message)
				}
				outputs = event.SerialCallEnded.Outputs
				return
			case nil != event.SerialLoopCallEnded && event.SerialLoopCallEnded.CallID == id:
				if nil != event.SerialLoopCallEnded.Error {
					err = errors.New(event.SerialLoopCallEnded.Error.Message)
				}
				outputs = event.SerialLoopCallEnded.Outputs
				return
			case nil != event.ParallelLoopCallEnded && event.ParallelLoopCallEnded.CallID == id:
				if nil != event.ParallelLoopCallEnded.Error {
					err = errors.New(event.ParallelLoopCallEnded.Error.Message)
				}
				outputs = event.ParallelLoopCallEnded.Outputs
				return
			case nil != event.ParallelCallEnded && event.ParallelCallEnded.CallID == id:
				if nil != event.ParallelCallEnded.Error {
					err = errors.New(event.ParallelCallEnded.Error.Message)
				}
				outputs = event.ParallelCallEnded.Outputs
				return
			case nil != event.OpKillRequested:
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
			rootOpID,
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
			rootOpID,
		)
	case nil != callSpec.Serial:
		clr.serialCaller.Call(
			callCtx,
			id,
			scope,
			rootOpID,
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
			rootOpID,
		)
	default:
		err = fmt.Errorf("Invalid call graph %+v\n", callSpec)
	}

}

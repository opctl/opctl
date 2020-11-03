package core

import (
	"context"
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
	) (
		map[string]*model.Value,
		error,
	)
}

func newCaller(
	callInterpreter call.Interpreter,
	containerCaller containerCaller,
	dataDirPath string,
	stateStore stateStore,
	pubSub pubsub.PubSub,
) caller {
	instance := &_caller{
		callInterpreter: callInterpreter,
		containerCaller: containerCaller,
		pubSub:          pubSub,
	}

	instance.opCaller = newOpCaller(
		stateStore,
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
) (
	map[string]*model.Value,
	error,
) {
	callCtx, cancelCall := context.WithCancel(ctx)
	defer cancelCall()
	var err error
	var isKilled bool
	var outputs map[string]*model.Value
	var call *model.Call
	callStartTime := time.Now().UTC()

	if nil != callCtx.Err() {
		// if context done NOOP
		return nil, nil
	}

	defer func() {
		// defer must be defined before conditional return statements so it always runs
		var outcome string
		if isKilled || nil != ctx.Err() {
			// this call or parent call killed/cancelled
			outcome = model.OpOutcomeKilled
		} else if nil != err {
			outcome = model.OpOutcomeFailed
		} else {
			outcome = model.OpOutcomeSucceeded
		}

		if nil == call {
			call = &model.Call{
				ID:     id,
				RootID: rootCallID,
			}
		}

		event := model.Event{
			CallEnded: &model.CallEnded{
				Call:    *call,
				Outcome: outcome,
				Outputs: outputs,
				Ref:     opPath,
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
		// NOOP
		return outputs, err
	}

	call, err = clr.callInterpreter.Interpret(
		scope,
		callSpec,
		id,
		opPath,
		parentCallID,
		rootCallID,
	)
	if nil != err {
		return nil, err
	}

	if nil != call.If && !*call.If {
		return outputs, err
	}

	clr.pubSub.Publish(
		model.Event{
			Timestamp: callStartTime,
			CallStarted: &model.CallStarted{
				Call: *call,
				Ref:  opPath,
			},
		},
	)

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
			case nil != event.CallKillRequested && event.CallKillRequested.Request.OpID == id:
				isKilled = true
				return
			}
		}
	}()

	switch {
	case nil != callSpec.Container:
		outputs, err = clr.containerCaller.Call(
			callCtx,
			call.Container,
			scope,
			callSpec.Container,
			rootCallID,
		)
	case nil != callSpec.Op:
		outputs, err = clr.opCaller.Call(
			callCtx,
			call.Op,
			scope,
			parentCallID,
			rootCallID,
			callSpec.Op,
		)
	case nil != callSpec.Parallel:
		outputs, err = clr.parallelCaller.Call(
			callCtx,
			id,
			scope,
			rootCallID,
			opPath,
			*callSpec.Parallel,
		)
	case nil != callSpec.ParallelLoop:
		outputs, err = clr.parallelLoopCaller.Call(
			callCtx,
			id,
			scope,
			*callSpec.ParallelLoop,
			opPath,
			parentCallID,
			rootCallID,
		)
	case nil != callSpec.Serial:
		outputs, err = clr.serialCaller.Call(
			callCtx,
			id,
			scope,
			rootCallID,
			opPath,
			*callSpec.Serial,
		)
	case nil != callSpec.SerialLoop:
		outputs, err = clr.serialLoopCaller.Call(
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

	return outputs, err
}

package core

import (
	"context"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/opctl/opctl/sdks/go/model"
	callpkg "github.com/opctl/opctl/sdks/go/opspec/interpreter/call"
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
	containerCaller containerCaller,
	dataDirPath string,
	pubSub pubsub.PubSub,
) caller {
	instance := &_caller{
		containerCaller: containerCaller,
		dataDirPath:     dataDirPath,
		pubSub:          pubSub,
	}

	instance.opCaller = newOpCaller(instance, dataDirPath)
	instance.parallelCaller = newParallelCaller(instance)
	instance.parallelLoopCaller = newParallelLoopCaller(instance)
	instance.serialCaller = newSerialCaller(instance)
	instance.serialLoopCaller = newSerialLoopCaller(instance)

	return instance
}

type _caller struct {
	containerCaller    containerCaller
	dataDirPath        string
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

	if callCtx.Err() != nil {
		// if context done NOOP
		return nil, nil
	}

	callStartTime := time.Now().UTC()

	// These variables are set by callers at the end of this function, but used
	// inside deferred func that emits call end events so they must be declared
	// early
	var err error
	var outputs map[string]*model.Value
	var call *model.Call
	var isKilled bool

	// emit a call ended event after this call is complete
	defer func() {
		// defer must be defined before conditional return statements so it always runs

		if call == nil {
			call = &model.Call{
				ID:       id,
				RootID:   rootCallID,
				ParentID: parentCallID,
			}
		}

		event := model.Event{
			CallEnded: &model.CallEnded{
				Call:    *call,
				Outputs: outputs,
				Ref:     opPath,
			},
			Timestamp: time.Now().UTC(),
		}

		if isKilled || ctx.Err() != nil {
			// this call or parent call killed/cancelled
			event.CallEnded.Outcome = model.OpOutcomeKilled
		} else if err != nil {
			event.CallEnded.Outcome = model.OpOutcomeFailed
			event.CallEnded.Error = &model.CallEndedError{
				Message: err.Error(),
			}
		} else {
			event.CallEnded.Outcome = model.OpOutcomeSucceeded
		}

		clr.pubSub.Publish(event)
	}()

	if callSpec == nil {
		// NOOP
		return outputs, err
	}

	call, err = callpkg.Interpret(
		ctx,
		scope,
		callSpec,
		id,
		opPath,
		parentCallID,
		rootCallID,
		clr.dataDirPath,
	)
	if err != nil {
		return nil, err
	}

	if call.If != nil && !*call.If {
		return outputs, err
	}

	// Ensure this is emitted just after the deferred operation to emit the end
	// event is set up, so we always have a matching start and end event
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
			case event.CallKillRequested != nil && event.CallKillRequested.Request.OpID == id:
				isKilled = true
				return
			}
		}
	}()

	switch {
	case callSpec.Container != nil:
		outputs, err = clr.containerCaller.Call(
			callCtx,
			call.Container,
			scope,
			callSpec.Container,
			rootCallID,
		)
	case callSpec.Op != nil:
		outputs, err = clr.opCaller.Call(
			callCtx,
			call.Op,
			scope,
			parentCallID,
			rootCallID,
			callSpec.Op,
		)
	case callSpec.Parallel != nil:
		outputs, err = clr.parallelCaller.Call(
			callCtx,
			id,
			scope,
			rootCallID,
			opPath,
			*callSpec.Parallel,
		)
	case callSpec.ParallelLoop != nil:
		outputs, err = clr.parallelLoopCaller.Call(
			callCtx,
			id,
			scope,
			*callSpec.ParallelLoop,
			opPath,
			parentCallID,
			rootCallID,
		)
	case callSpec.Serial != nil:
		outputs, err = clr.serialCaller.Call(
			callCtx,
			id,
			scope,
			rootCallID,
			opPath,
			*callSpec.Serial,
		)
	case callSpec.SerialLoop != nil:
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
		err = fmt.Errorf("invalid call graph '%+v'", callSpec)
	}

	return outputs, err
}

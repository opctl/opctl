package node

import (
	"context"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/pubsub"
	callpkg "github.com/opctl/opctl/sdks/go/opspec/interpreter/call"
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
	stateStore stateStore,
) caller {
	instance := &_caller{
		containerCaller: containerCaller,
		dataDirPath:     dataDirPath,
		pubSub:          pubSub,
		stateStore:      stateStore,
	}

	instance.opCaller = newOpCaller(
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
	containerCaller    containerCaller
	dataDirPath        string
	opCaller           opCaller
	parallelCaller     parallelCaller
	parallelLoopCaller parallelLoopCaller
	pubSub             pubsub.PubSub
	serialCaller       serialCaller
	serialLoopCaller   serialLoopCaller
	stateStore         stateStore
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

	if callCtx.Err() != nil {
		// if context done NOOP
		return nil, nil
	}

	defer func() {
		// defer must be defined before conditional return statements so it always runs

		if call == nil {
			call = &model.Call{
				ID:     id,
				RootID: rootCallID,
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

	if callSpec.Op != nil && callSpec.Op.PullCreds == nil {
		if auth := clr.stateStore.TryGetAuth(callSpec.Op.Ref); auth != nil {
			callSpec.Op.PullCreds = &model.CredsSpec{
				Username: auth.Username,
				Password: auth.Password,
			}
		}
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
			if panic := recover(); panic != nil {
				// recover from panics; treat as errors
				err = fmt.Errorf(
					"recovered from panic: %s\n%s", panic, string(debug.Stack()),
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

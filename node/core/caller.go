package core

//go:generate counterfeiter -o ./fakeCaller.go --fake-name fakeCaller ./ caller

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/opctl/sdk-golang/model"
	"github.com/opctl/sdk-golang/opspec/interpreter/call"
	"github.com/opctl/sdk-golang/util/pubsub"
)

type caller interface {
	// Call executes a call
	Call(
		ctx context.Context,
		id string,
		scope map[string]*model.Value,
		scg *model.SCG,
		opHandle model.DataHandle,
		parentCallID *string,
		rootOpID string,
	)
}

func newCaller(
	callInterpreter call.Interpreter,
	containerCaller containerCaller,
	dataDirPath string,
	callStore callStore,
	callKiller callKiller,
	pubSub pubsub.PubSub,
) caller {
	instance := &_caller{
		callInterpreter: callInterpreter,
		containerCaller: containerCaller,
		callStore:       callStore,
		pubSub:          pubSub,
	}

	instance.looper = newLooper(
		instance,
		pubSub,
	)

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

	instance.serialCaller = newSerialCaller(
		instance,
		pubSub,
	)

	return instance
}

type _caller struct {
	callInterpreter call.Interpreter
	containerCaller containerCaller
	callStore       callStore
	looper          looper
	opCaller        opCaller
	parallelCaller  parallelCaller
	pubSub          pubsub.PubSub
	serialCaller    serialCaller
}

func (clr _caller) Call(
	ctx context.Context,
	id string,
	scope map[string]*model.Value,
	scg *model.SCG,
	opHandle model.DataHandle,
	parentCallID *string,
	rootOpID string,
) {
	ctx, cancel := context.WithCancel(ctx)
	var err error
	var outputs map[string]*model.Value
	callStartTime := time.Now().UTC()

	defer func() {
		<-ctx.Done()

		// defer must be defined before conditional return statements so it always runs
		event := model.Event{
			CallEnded: &model.CallEndedEvent{
				CallID:     id,
				Outputs:    outputs,
				RootCallID: rootOpID,
			},
			Timestamp: time.Now().UTC(),
		}

		if nil != err {
			event.CallEnded.Error = &model.CallEndedEventError{
				Message: err.Error(),
			}
		}

		clr.pubSub.Publish(event)
	}()

	if nil == scg {
		cancel()
		// No Op
		return
	}

	go func() {
		defer cancel()

		eventChannel, _ := clr.pubSub.Subscribe(
			ctx,
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
			case nil != event.ParallelCallEnded && event.ParallelCallEnded.CallID == id:
				if nil != event.ParallelCallEnded.Error {
					err = errors.New(event.ParallelCallEnded.Error.Message)
				}
				outputs = event.ParallelCallEnded.Outputs
				return
			// if call killed, propogate to context
			case nil != event.CallKilled && event.CallKilled.CallID == id:
				return
			}
		}
	}()

	if nil != scg.Loop {
		clr.looper.Loop(
			ctx,
			id,
			scope,
			scg,
			opHandle,
			parentCallID,
			rootOpID,
		)
		return
	}

	var dcg *model.DCG
	dcg, err = clr.callInterpreter.Interpret(
		scope,
		scg,
		id,
		opHandle,
		parentCallID,
		rootOpID,
	)
	if nil != err {
		cancel()
		return
	}

	if nil != dcg.If && !*dcg.If {
		cancel()
		return
	}

	clr.callStore.Add(dcg)

	switch {
	case nil != scg.Container:
		clr.containerCaller.Call(
			ctx,
			dcg.Container,
			scope,
			scg.Container,
		)
		return
	case nil != scg.Op:
		clr.opCaller.Call(
			ctx,
			dcg.Op,
			scope,
			parentCallID,
			scg.Op,
		)
		return
	case len(scg.Parallel) > 0:
		clr.parallelCaller.Call(
			ctx,
			id,
			scope,
			rootOpID,
			opHandle,
			scg.Parallel,
		)
		return
	case len(scg.Serial) > 0:
		clr.serialCaller.Call(
			ctx,
			id,
			scope,
			rootOpID,
			opHandle,
			scg.Serial,
		)
		return
	default:
		err = fmt.Errorf("Invalid call graph %+v\n", scg)
		cancel()
		return
	}

}

package core

//go:generate counterfeiter -o ./fakeCaller.go --fake-name fakeCaller ./ caller

import (
	"context"
	"fmt"
	"sync"
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
	) error
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
		callKiller,
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
) error {
	ctx, cancel := context.WithCancel(ctx)
	var err error
	var outputs map[string]*model.Value
	var wg sync.WaitGroup

	defer func() {
		wg.Wait()
		cancel()

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
		// No Op
		return nil
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
		return err
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		eventChannel, _ := clr.pubSub.Subscribe(
			ctx,
			model.EventFilter{Roots: []string{rootOpID}},
		)

	eventLoop:
		for event := range eventChannel {
			switch {
			case nil != event.CallEnded && event.CallEnded.CallID == id:
				outputs = event.CallEnded.Outputs
				break eventLoop
			case nil != event.ContainerExited && event.ContainerExited.ContainerID == id:
				outputs = event.ContainerExited.Outputs
				break eventLoop
			case nil != event.OpEnded && event.OpEnded.OpID == id:
				outputs = event.OpEnded.Outputs
				break eventLoop
			case nil != event.SerialCallEnded && event.SerialCallEnded.CallID == id:
				outputs = event.SerialCallEnded.Outputs
				break eventLoop
			case nil != event.ParallelCallEnded && event.ParallelCallEnded.CallID == id:
				break eventLoop
			// if call killed, propogate to context
			case nil != event.CallKilled && event.CallKilled.CallID == id:
				cancel()
			}
		}
	}()

	if nil != dcg.Loop {
		err = clr.looper.Loop(
			ctx,
			id,
			scope,
			scg,
			opHandle,
			parentCallID,
			rootOpID,
		)
		return err
	}

	if nil != dcg.If && !*dcg.If {
		cancel()
		return nil
	}

	clr.callStore.Add(dcg)

	switch {
	case nil != scg.Container:
		err = clr.containerCaller.Call(
			ctx,
			dcg.Container,
			scope,
			scg.Container,
		)
		return err
	case nil != scg.Op:
		err = clr.opCaller.Call(
			ctx,
			dcg.Op,
			scope,
			parentCallID,
			scg.Op,
		)
		return err
	case len(scg.Parallel) > 0:
		err = clr.parallelCaller.Call(
			ctx,
			id,
			scope,
			rootOpID,
			opHandle,
			scg.Parallel,
		)
		return err
	case len(scg.Serial) > 0:
		err = clr.serialCaller.Call(
			ctx,
			id,
			scope,
			rootOpID,
			opHandle,
			scg.Serial,
		)
		return err
	default:
		err = fmt.Errorf("Invalid call graph %+v\n", scg)
		return err
	}

}

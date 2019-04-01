package core

//go:generate counterfeiter -o ./fakeCaller.go --fake-name fakeCaller ./ caller

import (
	"fmt"
	"time"

	"github.com/opctl/sdk-golang/model"
	"github.com/opctl/sdk-golang/opspec/interpreter/call"
	"github.com/opctl/sdk-golang/util/pubsub"
)

type caller interface {
	// Call executes a call
	Call(
		id string,
		scope map[string]*model.Value,
		scg *model.SCG,
		opHandle model.DataHandle,
		rootOpID string,
	) error
}

func newCaller(
	callInterpreter call.Interpreter,
	containerCaller containerCaller,
	dataDirPath string,
	dcgNodeRepo dcgNodeRepo,
	opKiller opKiller,
	pubSub pubsub.PubSub,
) caller {
	instance := &_caller{
		callInterpreter: callInterpreter,
		containerCaller: containerCaller,
		pubSub:          pubSub,
	}

	instance.looper = newLooper(
		instance,
		pubSub,
	)

	instance.opCaller = newOpCaller(
		pubSub,
		dcgNodeRepo,
		instance,
		dataDirPath,
	)

	instance.parallelCaller = newParallelCaller(
		instance,
		opKiller,
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
	looper          looper
	opCaller        opCaller
	parallelCaller  parallelCaller
	pubSub          pubsub.PubSub
	serialCaller    serialCaller
}

func (clr _caller) Call(
	id string,
	scope map[string]*model.Value,
	scg *model.SCG,
	opHandle model.DataHandle,
	rootOpID string,
) error {
	if nil == scg {
		// No Op
		return nil
	}

	dcg, err := clr.callInterpreter.Interpret(
		scope,
		scg,
		id,
		opHandle,
		rootOpID,
	)
	if nil != err {
		return err
	}

	defer func() {
		clr.pubSub.Publish(
			model.Event{
				CallEnded: &model.CallEndedEvent{
					CallID:     id,
					RootCallID: rootOpID,
				},
				Timestamp: time.Now().UTC(),
			},
		)
	}()

	if nil != dcg.Loop {
		return clr.looper.Loop(
			id,
			scope,
			scg,
			opHandle,
			rootOpID,
		)
	}

	if nil != dcg.If && !*dcg.If {
		return nil
	}

	switch {
	case nil != scg.Container:
		return clr.containerCaller.Call(
			dcg.Container,
			scope,
			scg.Container,
		)
	case nil != scg.Op:
		return clr.opCaller.Call(
			dcg.Op,
			scope,
			scg.Op,
		)
	case len(scg.Parallel) > 0:
		return clr.parallelCaller.Call(
			id,
			scope,
			rootOpID,
			opHandle,
			scg.Parallel,
		)
	case len(scg.Serial) > 0:
		return clr.serialCaller.Call(
			id,
			scope,
			rootOpID,
			opHandle,
			scg.Serial,
		)
	default:
		return fmt.Errorf("Invalid call graph %+v\n", scg)
	}

}

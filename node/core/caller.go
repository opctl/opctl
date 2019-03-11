package core

//go:generate counterfeiter -o ./fakeCaller.go --fake-name fakeCaller ./ caller

import (
	"fmt"

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
	pubSub pubsub.PubSub,
) *_caller {
	return &_caller{
		callInterpreter: callInterpreter,
		containerCaller: containerCaller,
		pubSub:          pubSub,
	}
}

type _caller struct {
	callInterpreter call.Interpreter
	containerCaller containerCaller
	opCaller        opCaller
	parallelCaller  parallelCaller
	pubSub          pubsub.PubSub
	serialCaller    serialCaller
}

func (this _caller) Call(
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

	dcg, err := this.callInterpreter.Interpret(
		scope,
		scg,
		id,
		opHandle,
		rootOpID,
	)
	if nil != err {
		return err
	}

	if nil != dcg.If && !*dcg.If {
		this.pubSub.Publish(
			model.Event{
				CallSkipped: &model.CallSkippedEvent{
					CallID:     id,
					RootCallID: rootOpID,
				},
			},
		)
		return nil
	}

	switch {
	case nil != scg.Container:
		return this.containerCaller.Call(
			dcg.Container,
			scope,
			id,
			scg.Container,
			opHandle,
			rootOpID,
		)
	case nil != scg.Op:
		return this.opCaller.Call(
			dcg.Op,
			scope,
			id,
			opHandle,
			rootOpID,
			scg.Op,
		)
	case len(scg.Parallel) > 0:
		return this.parallelCaller.Call(
			id,
			scope,
			rootOpID,
			opHandle,
			scg.Parallel,
		)
	case len(scg.Serial) > 0:
		return this.serialCaller.Call(
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

func (this *_caller) setOpCaller(
	opCaller opCaller,
) {
	this.opCaller = opCaller
}

func (this *_caller) setParallelCaller(
	parallelCaller parallelCaller,
) {
	this.parallelCaller = parallelCaller
}

func (this *_caller) setSerialCaller(
	serialCaller serialCaller,
) {
	this.serialCaller = serialCaller
}

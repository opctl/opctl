package core

import (
	"context"

	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/loop"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/loop/iteration"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/serialloop"

	"github.com/opctl/opctl/sdks/go/internal/uniquestring"
	"github.com/opctl/opctl/sdks/go/model"
)

//counterfeiter:generate -o internal/fakes/serialLoopCaller.go . serialLoopCaller
type serialLoopCaller interface {
	// Executes a serial loop call
	Call(
		ctx context.Context,
		id string,
		inboundScope map[string]*model.Value,
		callSpecSerialLoop model.SerialLoopCallSpec,
		opPath string,
		parentCallID *string,
		rootCallID string,
	) (
		map[string]*model.Value,
		error,
	)
}

func newSerialLoopCaller(caller caller) serialLoopCaller {
	return _serialLoopCaller{
		caller: caller,
	}
}

type _serialLoopCaller struct {
	caller caller
}

func (lpr _serialLoopCaller) Call(
	ctx context.Context,
	id string,
	inboundScope map[string]*model.Value,
	callSpecSerialLoop model.SerialLoopCallSpec,
	opPath string,
	parentCallID *string,
	rootCallID string,
) (
	map[string]*model.Value,
	error,
) {
	index := 0
	scope, err := iteration.Scope(
		index,
		inboundScope,
		callSpecSerialLoop.Range,
		callSpecSerialLoop.Vars,
	)
	if err != nil {
		return nil, err
	}

	// interpret initial iteration of the loop
	callSerialLoop, err := serialloop.Interpret(
		callSpecSerialLoop,
		scope,
	)
	if err != nil {
		return nil, err
	}

	for !serialloop.IsIterationComplete(index, callSerialLoop) {
		callID, err := uniquestring.Construct()
		if err != nil {
			return nil, err
		}

		outputs, err := lpr.caller.Call(
			ctx,
			callID,
			scope,
			&callSpecSerialLoop.Run,
			opPath,
			parentCallID,
			rootCallID,
		)
		if err != nil {
			return nil, err
		}

		for name, value := range outputs {
			scope[name] = value
		}

		index++

		if serialloop.IsIterationComplete(index, callSerialLoop) {
			break
		}

		scope, err = iteration.Scope(
			index,
			scope,
			callSpecSerialLoop.Range,
			callSpecSerialLoop.Vars,
		)
		if err != nil {
			return nil, err
		}

		// interpret next iteration of the loop
		callSerialLoop, err = serialloop.Interpret(
			callSpecSerialLoop,
			scope,
		)
		if err != nil {
			return nil, err
		}
	}

	outboundScope := loop.DeScope(
		inboundScope,
		callSpecSerialLoop.Range,
		callSpecSerialLoop.Vars,
		scope,
	)

	return outboundScope, err
}

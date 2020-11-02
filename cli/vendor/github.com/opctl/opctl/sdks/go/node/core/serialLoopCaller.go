package core

import (
	"context"
	"errors"
	"time"

	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/loop"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/loop/iteration"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/serialloop"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/loopable"

	"github.com/opctl/opctl/sdks/go/internal/uniquestring"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/pubsub"
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

func newSerialLoopCaller(
	caller caller,
	pubSub pubsub.PubSub,
) serialLoopCaller {
	return _serialLoopCaller{
		caller:                caller,
		iterationScoper:       iteration.NewScoper(),
		loopableInterpreter:   loopable.NewInterpreter(),
		loopDeScoper:          loop.NewDeScoper(),
		pubSub:                pubSub,
		serialLoopInterpreter: serialloop.NewInterpreter(),
		uniqueStringFactory:   uniquestring.NewUniqueStringFactory(),
	}
}

type _serialLoopCaller struct {
	caller                caller
	iterationScoper       iteration.Scoper
	loopableInterpreter   loopable.Interpreter
	loopDeScoper          loop.DeScoper
	pubSub                pubsub.PubSub
	serialLoopInterpreter serialloop.Interpreter
	uniqueStringFactory   uniquestring.UniqueStringFactory
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
	var err error
	outboundScope := map[string]*model.Value{}
	var callSerialLoop *model.SerialLoopCall

	index := 0
	outboundScope, err = lpr.iterationScoper.Scope(
		index,
		inboundScope,
		callSpecSerialLoop.Range,
		callSpecSerialLoop.Vars,
	)
	if nil != err {
		return outboundScope, err
	}

	// interpret initial iteration of the loop
	callSerialLoop, err = lpr.serialLoopInterpreter.Interpret(
		callSpecSerialLoop,
		outboundScope,
	)
	if nil != err {
		return outboundScope, err
	}

	for !serialloop.IsIterationComplete(index, callSerialLoop) {
		eventFilterSince := time.Now().UTC()

		var callID string
		callID, err = lpr.uniqueStringFactory.Construct()
		if nil != err {
			return outboundScope, err
		}

		lpr.caller.Call(
			ctx,
			callID,
			outboundScope,
			&callSpecSerialLoop.Run,
			opPath,
			parentCallID,
			rootCallID,
		)

		// subscribe to events
		// @TODO: handle err channel
		eventChannel, _ := lpr.pubSub.Subscribe(
			ctx,
			model.EventFilter{
				Roots: []string{rootCallID},
				Since: &eventFilterSince,
			},
		)

	eventLoop:
		for event := range eventChannel {
			// merge child outboundScope w/ outboundScope, child outboundScope having precedence
			switch {
			case nil != event.CallEnded && event.CallEnded.Call.Id == callID:
				if nil != event.CallEnded.Error {
					err = errors.New(event.CallEnded.Error.Message)
					return outboundScope, err
				}
				for name, value := range event.CallEnded.Outputs {
					outboundScope[name] = value
				}
				break eventLoop
			}
		}

		index++

		if serialloop.IsIterationComplete(index, callSerialLoop) {
			break
		}

		outboundScope, err = lpr.iterationScoper.Scope(
			index,
			outboundScope,
			callSpecSerialLoop.Range,
			callSpecSerialLoop.Vars,
		)
		if nil != err {
			return outboundScope, err
		}

		// interpret next iteration of the loop
		callSerialLoop, err = lpr.serialLoopInterpreter.Interpret(
			callSpecSerialLoop,
			outboundScope,
		)
		if nil != err {
			return outboundScope, err
		}
	}

	outboundScope = lpr.loopDeScoper.DeScope(
		inboundScope,
		callSpecSerialLoop.Range,
		callSpecSerialLoop.Vars,
		outboundScope,
	)

	return outboundScope, err
}

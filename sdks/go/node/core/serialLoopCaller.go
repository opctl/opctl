package core

//go:generate counterfeiter -o ./fakeSerialLoopCaller.go --fake-name fakeSerialLoopCaller ./ serialLoopCaller

import (
	"context"
	"errors"
	"time"

	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/loop"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/loop/iteration"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/serialloop"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/loopable"

	"github.com/opctl/opctl/sdks/go/types"
	"github.com/opctl/opctl/sdks/go/util/pubsub"
	"github.com/opctl/opctl/sdks/go/util/uniquestring"
)

type serialLoopCaller interface {
	// Executes a serial loop call
	Call(
		ctx context.Context,
		id string,
		inboundScope map[string]*types.Value,
		scgSerialLoop types.SCGSerialLoopCall,
		opHandle types.DataHandle,
		parentCallID *string,
		rootOpID string,
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
	inboundScope map[string]*types.Value,
	scgSerialLoop types.SCGSerialLoopCall,
	opHandle types.DataHandle,
	parentCallID *string,
	rootOpID string,
) {
	var err error
	outboundScope := map[string]*types.Value{}

	defer func() {
		// defer must be defined before conditional return statements so it always runs
		event := types.Event{
			Timestamp: time.Now().UTC(),
			SerialLoopCallEnded: &types.SerialLoopCallEndedEvent{
				CallID:   id,
				RootOpID: rootOpID,
				Outputs:  outboundScope,
			},
		}

		if nil != err {
			event.SerialLoopCallEnded.Error = &types.CallEndedEventError{
				Message: err.Error(),
			}
		}

		lpr.pubSub.Publish(event)
	}()

	index := 0
	outboundScope, err = lpr.iterationScoper.Scope(
		index,
		inboundScope,
		scgSerialLoop.Range,
		scgSerialLoop.Vars,
		opHandle,
	)
	if nil != err {
		return
	}

	// interpret initial iteration of the loop
	var dcgSerialLoop *types.DCGSerialLoopCall
	dcgSerialLoop, err = lpr.serialLoopInterpreter.Interpret(
		opHandle,
		scgSerialLoop,
		outboundScope,
	)
	if nil != err {
		return
	}

	for !serialloop.IsIterationComplete(index, dcgSerialLoop) {
		eventFilterSince := time.Now().UTC()

		var callID string
		callID, err = lpr.uniqueStringFactory.Construct()
		if nil != err {
			return
		}

		lpr.caller.Call(
			ctx,
			callID,
			outboundScope,
			&scgSerialLoop.Run,
			opHandle,
			parentCallID,
			rootOpID,
		)

		// subscribe to events
		// @TODO: handle err channel
		eventChannel, _ := lpr.pubSub.Subscribe(
			ctx,
			types.EventFilter{
				Roots: []string{rootOpID},
				Since: &eventFilterSince,
			},
		)

	eventLoop:
		for event := range eventChannel {
			// merge child outboundScope w/ outboundScope, child outboundScope having precedence
			switch {
			case nil != event.CallEnded && event.CallEnded.CallID == callID:
				if nil != event.CallEnded.Error {
					err = errors.New(event.CallEnded.Error.Message)
					return
				}
				for name, value := range event.CallEnded.Outputs {
					outboundScope[name] = value
				}
				break eventLoop
			}
		}

		index++

		if serialloop.IsIterationComplete(index, dcgSerialLoop) {
			break
		}

		outboundScope, err = lpr.iterationScoper.Scope(
			index,
			outboundScope,
			scgSerialLoop.Range,
			scgSerialLoop.Vars,
			opHandle,
		)
		if nil != err {
			return
		}

		// interpret next iteration of the loop
		dcgSerialLoop, err = lpr.serialLoopInterpreter.Interpret(
			opHandle,
			scgSerialLoop,
			outboundScope,
		)
		if nil != err {
			return
		}
	}

	outboundScope = lpr.loopDeScoper.DeScope(
		inboundScope,
		scgSerialLoop.Range,
		scgSerialLoop.Vars,
		outboundScope,
	)
}

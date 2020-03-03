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
		scgSerialLoop model.SCGSerialLoopCall,
		opHandle model.DataHandle,
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
	inboundScope map[string]*model.Value,
	scgSerialLoop model.SCGSerialLoopCall,
	opHandle model.DataHandle,
	parentCallID *string,
	rootOpID string,
) {
	var err error
	outboundScope := map[string]*model.Value{}

	defer func() {
		// defer must be defined before conditional return statements so it always runs
		event := model.Event{
			Timestamp: time.Now().UTC(),
			SerialLoopCallEnded: &model.SerialLoopCallEndedEvent{
				CallID:   id,
				RootOpID: rootOpID,
				Outputs:  outboundScope,
			},
		}

		if nil != err {
			event.SerialLoopCallEnded.Error = &model.CallEndedEventError{
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
	var dcgSerialLoop *model.DCGSerialLoopCall
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
			model.EventFilter{
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

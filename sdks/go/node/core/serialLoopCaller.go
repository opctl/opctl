package core

import (
	"context"
	"errors"
	"time"

	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/loop"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/loop/iteration"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/serialloop"

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
		caller: caller,
		pubSub: pubSub,
	}
}

type _serialLoopCaller struct {
	caller caller
	pubSub pubsub.PubSub
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
	outboundScope := map[string]*model.Value{}
	var callSerialLoop *model.SerialLoopCall

	index := 0
	outboundScope, err := iteration.Scope(
		index,
		inboundScope,
		callSpecSerialLoop.Range,
		callSpecSerialLoop.Vars,
	)
	if err != nil {
		return nil, err
	}

	// interpret initial iteration of the loop
	callSerialLoop, err = serialloop.Interpret(
		callSpecSerialLoop,
		outboundScope,
	)
	if err != nil {
		return nil, err
	}

	for !serialloop.IsIterationComplete(index, callSerialLoop) {
		eventFilterSince := time.Now().UTC()

		var callID string
		callID, err = uniquestring.Construct()
		if err != nil {
			return nil, err
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
			case event.CallEnded != nil && event.CallEnded.Call.ID == callID:
				if event.CallEnded.Error != nil {
					err = errors.New(event.CallEnded.Error.Message)
					return nil, err
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

		outboundScope, err = iteration.Scope(
			index,
			outboundScope,
			callSpecSerialLoop.Range,
			callSpecSerialLoop.Vars,
		)
		if err != nil {
			return nil, err
		}

		// interpret next iteration of the loop
		callSerialLoop, err = serialloop.Interpret(
			callSpecSerialLoop,
			outboundScope,
		)
		if err != nil {
			return nil, err
		}
	}

	outboundScope = loop.DeScope(
		inboundScope,
		callSpecSerialLoop.Range,
		callSpecSerialLoop.Vars,
		outboundScope,
	)

	return outboundScope, err
}

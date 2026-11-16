package node

import (
	"context"
	"errors"
	"time"

	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/loop"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/loop/iteration"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/serialloop"

	"github.com/opctl/opctl/sdks/go/internal/uniquestring"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/pubsub"
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

	// subscribe to events
	// @TODO: handle err channel
	subCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	eventFilterSince := time.Now().UTC()
	eventChannel, _ := lpr.pubSub.Subscribe(
		subCtx,
		model.EventFilter{
			Roots: []string{rootCallID},
			Since: &eventFilterSince,
		},
	)
	watchChildCallEnded := newChildCallEndedWatcher(eventChannel)

	for !serialloop.IsIterationComplete(index, callSerialLoop) {
		var callID string
		callID, err = uniquestring.Construct()
		if err != nil {
			return nil, err
		}

		childCallEndedChannel := watchChildCallEnded(callID)
		lpr.caller.Call(
			ctx,
			callID,
			outboundScope,
			&callSpecSerialLoop.Run,
			opPath,
			parentCallID,
			rootCallID,
		)

		if callEnded := <-childCallEndedChannel; callEnded != nil {
			if callEnded.Error != nil {
				err = errors.New(callEnded.Error.Message)
				return nil, err
			}
			// merge child outboundScope w/ outboundScope, child outboundScope having precedence
			for name, value := range callEnded.Outputs {
				outboundScope[name] = value
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

package core

import (
	"context"
	"errors"
	"time"

	"github.com/opctl/opctl/sdks/go/internal/uniquestring"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/pubsub"
)

//counterfeiter:generate -o internal/fakes/serialCaller.go . serialCaller
type serialCaller interface {
	// Executes a serial call
	Call(
		ctx context.Context,
		callID string,
		inboundScope map[string]*model.Value,
		rootCallID string,
		opPath string,
		callSpecSerialCall []*model.CallSpec,
	) (
		map[string]*model.Value,
		error,
	)
}

func newSerialCaller(
	caller caller,
	pubSub pubsub.PubSub,
) serialCaller {

	return _serialCaller{
		caller:              caller,
		pubSub:              pubSub,
		uniqueStringFactory: uniquestring.NewUniqueStringFactory(),
	}

}

type _serialCaller struct {
	caller              caller
	pubSub              pubsub.PubSub
	uniqueStringFactory uniquestring.UniqueStringFactory
}

func (sc _serialCaller) Call(
	ctx context.Context,
	callID string,
	inboundScope map[string]*model.Value,
	rootCallID string,
	opPath string,
	callSpecSerialCall []*model.CallSpec,
) (
	map[string]*model.Value,
	error,
) {
	outputs := map[string]*model.Value{}
	for varName, varData := range inboundScope {
		outputs[varName] = varData
	}

	// subscribe to events
	// @TODO: handle err channel
	eventFilterSince := time.Now().UTC()
	eventChannel, _ := sc.pubSub.Subscribe(
		ctx,
		model.EventFilter{
			Roots: []string{rootCallID},
			Since: &eventFilterSince,
		},
	)

	for _, callSpecCall := range callSpecSerialCall {

		var childCallID string
		childCallID, err := sc.uniqueStringFactory.Construct()
		if nil != err {
			// end run immediately on any error
			return nil, err
		}

		sc.caller.Call(
			ctx,
			childCallID,
			outputs,
			callSpecCall,
			opPath,
			&callID,
			rootCallID,
		)

	eventLoop:
		for event := range eventChannel {
			// merge child outboundScope w/ outboundScope, child outboundScope having precedence
			switch {
			case nil != event.CallEnded && event.CallEnded.Call.ID == childCallID:
				if nil != event.CallEnded.Error {
					// end on any error
					return nil, errors.New(event.CallEnded.Error.Message)
				}
				for name, value := range event.CallEnded.Outputs {
					outputs[name] = value
				}
				break eventLoop
			}
		}

	}

	return outputs, nil
}

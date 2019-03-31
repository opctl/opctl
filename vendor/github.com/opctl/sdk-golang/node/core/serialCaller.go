package core

//go:generate counterfeiter -o ./fakeSerialCaller.go --fake-name fakeSerialCaller ./ serialCaller

import (
	"context"
	"time"

	"github.com/opctl/sdk-golang/model"
	"github.com/opctl/sdk-golang/util/pubsub"
	"github.com/opctl/sdk-golang/util/uniquestring"
)

type serialCaller interface {
	// Executes a serial call
	Call(
		callId string,
		inboundScope map[string]*model.Value,
		rootOpID string,
		opHandle model.DataHandle,
		scgSerialCall []*model.SCG,
	) error
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

func (this _serialCaller) Call(
	callId string,
	inboundScope map[string]*model.Value,
	rootOpID string,
	opHandle model.DataHandle,
	scgSerialCall []*model.SCG,
) error {
	outputs := map[string]*model.Value{}
	for varName, varData := range inboundScope {
		outputs[varName] = varData
	}

	defer func() {
		// defer must be defined before conditional return statements so it always runs
		this.pubSub.Publish(
			model.Event{
				Timestamp: time.Now().UTC(),
				SerialCallEnded: &model.SerialCallEndedEvent{
					CallID:   callId,
					RootOpID: rootOpID,
					Outputs:  outputs,
				},
			},
		)
	}()

	ctx := context.TODO()

	for _, scgCall := range scgSerialCall {
		eventFilterSince := time.Now().UTC()

		childCallID, err := this.uniqueStringFactory.Construct()
		if nil != err {
			// end run immediately on any error
			return err
		}

		if err := this.caller.Call(
			childCallID,
			outputs,
			scgCall,
			opHandle,
			rootOpID,
		); nil != err {
			// end run immediately on any error
			return err
		}

		// subscribe to events
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()
		// @TODO: handle err channel
		eventChannel, _ := this.pubSub.Subscribe(
			ctx,
			model.EventFilter{
				Roots: []string{rootOpID},
				Since: &eventFilterSince,
			},
		)

	eventLoop:
		for event := range eventChannel {
			// merge child outputs w/ outputs, child outputs having precedence
			switch {
			case nil != event.OpEnded && event.OpEnded.OpID == childCallID:
				for name, value := range event.OpEnded.Outputs {
					outputs[name] = value
				}
				break eventLoop
			case nil != event.ContainerExited && event.ContainerExited.ContainerID == childCallID:
				for name, value := range event.ContainerExited.Outputs {
					outputs[name] = value
				}
				break eventLoop
			case nil != event.SerialCallEnded && event.SerialCallEnded.CallID == childCallID:
				for name, value := range event.SerialCallEnded.Outputs {
					outputs[name] = value
				}
				break eventLoop
			case nil != event.ParallelCallEnded && event.ParallelCallEnded.CallID == childCallID:
				break eventLoop
			case nil != event.CallEnded && event.CallEnded.CallID == childCallID:
				break eventLoop
			}
		}
		cancel()

	}

	return nil

}

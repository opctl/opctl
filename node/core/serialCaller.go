package core

//go:generate counterfeiter -o ./fakeSerialCaller.go --fake-name fakeSerialCaller ./ serialCaller

import (
	"context"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/util/pubsub"
	"github.com/opspec-io/sdk-golang/util/uniquestring"
	"time"
)

type serialCaller interface {
	// Executes a serial call
	Call(
		callId string,
		inboundScope map[string]*model.Value,
		rootOpId string,
		opDirHandle model.DataHandle,
		scgSerialCall []*model.SCG,
	) error
}

func newSerialCaller(
	caller caller,
	pubSub pubsub.PubSub,
	uniqueStringFactory uniquestring.UniqueStringFactory,
) serialCaller {

	return _serialCaller{
		caller:              caller,
		pubSub:              pubSub,
		uniqueStringFactory: uniqueStringFactory,
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
	rootOpId string,
	opDirHandle model.DataHandle,
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
					CallId:   callId,
					RootOpId: rootOpId,
					Outputs:  outputs,
				},
			},
		)
	}()

	ctx := context.TODO()

	for _, scgCall := range scgSerialCall {
		eventFilterSince := time.Now().UTC()

		childCallId, err := this.uniqueStringFactory.Construct()
		if nil != err {
			// end run immediately on any error
			return err
		}

		if err := this.caller.Call(
			childCallId,
			outputs,
			scgCall,
			opDirHandle,
			rootOpId,
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
				Roots: []string{rootOpId},
				Since: &eventFilterSince,
			},
		)

	eventLoop:
		for event := range eventChannel {
			// merge child outputs w/ outputs, child outputs having precedence
			switch {
			case nil != event.OpEnded && event.OpEnded.OpId == childCallId:
				for name, value := range event.OpEnded.Outputs {
					outputs[name] = value
				}
				break eventLoop
			case nil != event.ContainerExited && event.ContainerExited.ContainerId == childCallId:
				for name, value := range event.ContainerExited.Outputs {
					outputs[name] = value
				}
				break eventLoop
			case nil != event.SerialCallEnded && event.SerialCallEnded.CallId == childCallId:
				for name, value := range event.SerialCallEnded.Outputs {
					outputs[name] = value
				}
				break eventLoop
			case nil != event.ParallelCallEnded && event.ParallelCallEnded.CallId == childCallId:
				break eventLoop
			}
		}
		cancel()

	}

	return nil

}

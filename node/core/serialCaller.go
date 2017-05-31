package core

//go:generate counterfeiter -o ./fakeSerialCaller.go --fake-name fakeSerialCaller ./ serialCaller

import (
	"github.com/opctl/opctl/util/pubsub"
	"github.com/opctl/opctl/util/uniquestring"
	"github.com/opspec-io/sdk-golang/model"
	"time"
)

type serialCaller interface {
	// Executes a serial call
	Call(
		callId string,
		inboundScope map[string]*model.Data,
		rootOpId string,
		pkgRef string,
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
	inboundScope map[string]*model.Data,
	rootOpId string,
	pkgRef string,
	scgSerialCall []*model.SCG,
) error {
	outputs := map[string]*model.Data{}
	for varName, varData := range inboundScope {
		outputs[varName] = varData
	}

	defer func() {
		// defer must be defined before conditional return statements so it always runs
		this.pubSub.Publish(
			&model.Event{
				Timestamp: time.Now().UTC(),
				SerialCallEnded: &model.SerialCallEndedEvent{
					CallId:   callId,
					RootOpId: rootOpId,
					Outputs:  outputs,
				},
			},
		)
	}()

	for _, scgCall := range scgSerialCall {
		eventFilterSince := time.Now().UTC()
		childCallId := this.uniqueStringFactory.Construct()
		if err := this.caller.Call(
			childCallId,
			outputs,
			scgCall,
			pkgRef,
			rootOpId,
		); nil != err {
			// end run immediately on any error
			return err
		}

		// subscribe to events
		eventChannel := make(chan *model.Event, 150)
		this.pubSub.Subscribe(
			&model.EventFilter{
				RootOpIds: []string{rootOpId},
				Since:     &eventFilterSince,
			},
			eventChannel,
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

	}

	return nil

}

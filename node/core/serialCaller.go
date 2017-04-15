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
	) (
		err error,
	)
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
) (
	err error,
) {
	defer func() {
		// defer must be defined before conditional return statements so it always runs

		this.pubSub.Publish(
			&model.Event{
				Timestamp: time.Now().UTC(),
				SerialCallEnded: &model.SerialCallEndedEvent{
					CallId:   callId,
					RootOpId: rootOpId,
				},
			},
		)

	}()

	scope := map[string]*model.Data{}
	for varName, varData := range inboundScope {
		scope[varName] = varData
	}

	eventFilterSince := time.Now().UTC()
	for _, scgCall := range scgSerialCall {
		childCallId := this.uniqueStringFactory.Construct()
		err = this.caller.Call(
			childCallId,
			scope,
			scgCall,
			pkgRef,
			rootOpId,
		)
		if nil != err {
			// end run immediately on any error
			return
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

		// send outputs
	eventLoop:
		for event := range eventChannel {
			switch {
			case nil != event.OpEnded && event.OpEnded.OpId == childCallId:
				break eventLoop
			case nil != event.ContainerExited && event.ContainerExited.ContainerId == childCallId:
				break eventLoop
			case nil != event.SerialCallEnded && event.SerialCallEnded.CallId == childCallId:
				break eventLoop
			case nil != event.ParallelCallEnded && event.ParallelCallEnded.CallId == childCallId:
				break eventLoop
			case nil != event.OutputInitialized && event.OutputInitialized.CallId == childCallId:
				childOutput := event.OutputInitialized
				if scgOpCall := scgCall.Op; nil != scgCall.Op {
					// apply bound child outputs to current scope
					for currentScopeVarName, childScopeVarName := range scgOpCall.Outputs {
						if currentScopeVarName == childOutput.Name || childScopeVarName == childOutput.Name {
							scope[currentScopeVarName] = childOutput.Value
						}
					}
				} else {
					// apply child outputs to current scope
					scope[childOutput.Name] = childOutput.Value
				}
			}
		}

	}

	// @TODO: stream outputs from last child
	for childOutputName, childOutputValue := range scope {
		this.pubSub.Publish(&model.Event{
			OutputInitialized: &model.OutputInitializedEvent{
				Name:     childOutputName,
				Value:    childOutputValue,
				RootOpId: rootOpId,
				CallId:   callId,
			},
		})
	}

	return

}

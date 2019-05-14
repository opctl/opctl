package core

//go:generate counterfeiter -o ./fakeOpCaller.go --fake-name fakeOpCaller ./ opCaller

import (
	"context"
	"time"

	"github.com/opctl/sdk-golang/model"
	"github.com/opctl/sdk-golang/opspec/interpreter/call/op/outputs"
	dotyml "github.com/opctl/sdk-golang/opspec/opfile"
	"github.com/opctl/sdk-golang/util/pubsub"
)

type opCaller interface {
	// Executes an op call
	Call(
		dcgOpCall *model.DCGOpCall,
		inboundScope map[string]*model.Value,
		parentCallID *string,
		scgOpCall *model.SCGOpCall,
	) error
}

func newOpCaller(
	callStore callStore,
	pubSub pubsub.PubSub,
	caller caller,
	dataDirPath string,
) opCaller {
	return _opCaller{
		caller:             caller,
		callStore:          callStore,
		outputsInterpreter: outputs.NewInterpreter(),
		dotYmlGetter:       dotyml.NewGetter(),
		pubSub:             pubSub,
	}
}

type _opCaller struct {
	callStore          callStore
	caller             caller
	outputsInterpreter outputs.Interpreter
	dotYmlGetter       dotyml.Getter
	pubSub             pubsub.PubSub
}

func (oc _opCaller) Call(
	dcgOpCall *model.DCGOpCall,
	inboundScope map[string]*model.Value,
	parentCallID *string,
	scgOpCall *model.SCGOpCall,
) error {
	var err error
	outboundScope := map[string]*model.Value{}

	defer func() {
		// defer must be defined before conditional return statements so it always runs
		var opOutcome string
		if oc.callStore.Get(dcgOpCall.OpID).IsKilled {
			opOutcome = model.OpOutcomeKilled
		} else if nil != err {
			oc.pubSub.Publish(
				model.Event{
					Timestamp: time.Now().UTC(),
					OpErred: &model.OpErredEvent{
						Msg:      err.Error(),
						OpID:     dcgOpCall.OpID,
						OpRef:    dcgOpCall.OpHandle.Ref(),
						RootOpID: dcgOpCall.RootOpID,
					},
				},
			)
			opOutcome = model.OpOutcomeFailed
		} else {
			opOutcome = model.OpOutcomeSucceeded
		}

		oc.pubSub.Publish(
			model.Event{
				Timestamp: time.Now().UTC(),
				OpEnded: &model.OpEndedEvent{
					OpID:     dcgOpCall.OpID,
					OpRef:    dcgOpCall.OpHandle.Ref(),
					Outcome:  opOutcome,
					RootOpID: dcgOpCall.RootOpID,
					Outputs:  outboundScope,
				},
			},
		)

	}()

	oc.pubSub.Publish(
		model.Event{
			Timestamp: time.Now().UTC(),
			OpStarted: &model.OpStartedEvent{
				OpID:     dcgOpCall.OpID,
				OpRef:    dcgOpCall.OpHandle.Ref(),
				RootOpID: dcgOpCall.RootOpID,
			},
		},
	)

	// establish op output channel
	opOutputsChan := make(chan map[string]*model.Value, 1)
	go func() {
		opOutputsChan <- oc.waitOnOpOutputs(
			context.TODO(),
			dcgOpCall,
		)
	}()

	err = oc.caller.Call(
		dcgOpCall.ChildCallID,
		dcgOpCall.Inputs,
		dcgOpCall.ChildCallSCG,
		dcgOpCall.OpHandle,
		parentCallID,
		dcgOpCall.RootOpID,
	)

	if nil != err {
		return err
	}

	// wait on op outboundScope
	opOutputs := <-opOutputsChan

	opDotYml, err := oc.dotYmlGetter.Get(
		context.TODO(),
		dcgOpCall.OpHandle,
	)
	if nil != err {
		return err
	}
	opPath := dcgOpCall.OpHandle.Path()
	opOutputs, err = oc.outputsInterpreter.Interpret(opOutputs, opDotYml.Outputs, *opPath)

	// filter op outboundScope to bound call outboundScope
	for boundName, boundValue := range scgOpCall.Outputs {
		// return bound outboundScope
		if "" == boundValue {
			// implicit value
			boundValue = boundName
		}
		for opOutputName, opOutputValue := range opOutputs {
			if boundValue == opOutputName {
				outboundScope[boundName] = opOutputValue
			}
		}
	}

	return err
}

func (oc _opCaller) waitOnOpOutputs(
	ctx context.Context,
	dcgOpCall *model.DCGOpCall,
) map[string]*model.Value {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// subscribe to events
	eventFilterSince := time.Now().UTC()
	eventChannel, _ := oc.pubSub.Subscribe(
		ctx,
		model.EventFilter{
			Roots: []string{dcgOpCall.RootOpID},
			Since: &eventFilterSince,
		},
	)

	opOutputs := map[string]*model.Value{}
eventLoop:
	for event := range eventChannel {
		switch {
		case nil != event.OpEnded && event.OpEnded.OpID == dcgOpCall.OpID:
			// parent ended prematurely
			return nil
		case nil != event.OpEnded && event.OpEnded.OpID == dcgOpCall.ChildCallID:
			for name, value := range event.OpEnded.Outputs {
				opOutputs[name] = value
			}
			break eventLoop
		case nil != event.ContainerExited && event.ContainerExited.ContainerID == dcgOpCall.ChildCallID:
			for name, value := range event.ContainerExited.Outputs {
				opOutputs[name] = value
			}
			break eventLoop
		case nil != event.SerialCallEnded && event.SerialCallEnded.CallID == dcgOpCall.ChildCallID:
			for name, value := range event.SerialCallEnded.Outputs {
				opOutputs[name] = value
			}
			break eventLoop
		case nil != event.ParallelCallEnded && event.ParallelCallEnded.CallID == dcgOpCall.ChildCallID:
			// parallel calls have no outboundScope
			return nil
		case nil != event.CallEnded && event.CallEnded.CallID == dcgOpCall.ChildCallID:
			for name, value := range event.CallEnded.Outputs {
				opOutputs[name] = value
			}
			break eventLoop
		}
	}

	return opOutputs
}

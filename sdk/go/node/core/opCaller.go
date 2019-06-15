package core

//go:generate counterfeiter -o ./fakeOpCaller.go --fake-name fakeOpCaller ./ opCaller

import (
	"context"
	"errors"
	"time"

	"github.com/opctl/sdk-golang/model"
	"github.com/opctl/sdk-golang/opspec/interpreter/call/op/outputs"
	dotyml "github.com/opctl/sdk-golang/opspec/opfile"
	"github.com/opctl/sdk-golang/util/pubsub"
)

type opCaller interface {
	// Executes an op call
	Call(
		ctx context.Context,
		dcgOpCall *model.DCGOpCall,
		inboundScope map[string]*model.Value,
		parentCallID *string,
		scgOpCall *model.SCGOpCall,
	)
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
	ctx context.Context,
	dcgOpCall *model.DCGOpCall,
	inboundScope map[string]*model.Value,
	parentCallID *string,
	scgOpCall *model.SCGOpCall,
) {
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

	opStartedTime := time.Now().UTC()

	oc.pubSub.Publish(
		model.Event{
			Timestamp: opStartedTime,
			OpStarted: &model.OpStartedEvent{
				OpID:     dcgOpCall.OpID,
				OpRef:    dcgOpCall.OpHandle.Ref(),
				RootOpID: dcgOpCall.RootOpID,
			},
		},
	)

	oc.caller.Call(
		ctx,
		dcgOpCall.ChildCallID,
		dcgOpCall.Inputs,
		dcgOpCall.ChildCallSCG,
		dcgOpCall.OpHandle,
		&dcgOpCall.OpID,
		dcgOpCall.RootOpID,
	)

	// subscribe to events
	eventChannel, _ := oc.pubSub.Subscribe(
		ctx,
		model.EventFilter{
			Roots: []string{dcgOpCall.RootOpID},
			Since: &opStartedTime,
		},
	)

	opOutputs := map[string]*model.Value{}

eventLoop:
	for event := range eventChannel {
		switch {
		case nil != event.OpEnded && event.OpEnded.OpID == dcgOpCall.OpID:
			// parent ended prematurely
			return
		case nil != event.CallEnded && event.CallEnded.CallID == dcgOpCall.ChildCallID:
			if nil != event.CallEnded.Error {
				err = errors.New(event.CallEnded.Error.Message)
				// end on any error
				return
			}
			for name, value := range event.CallEnded.Outputs {
				opOutputs[name] = value
			}
			break eventLoop
		}
	}

	var opDotYml *model.OpDotYml
	opDotYml, err = oc.dotYmlGetter.Get(
		ctx,
		dcgOpCall.OpHandle,
	)
	if nil != err {
		return
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
}

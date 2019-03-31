package core

//go:generate counterfeiter -o ./fakeOpCaller.go --fake-name fakeOpCaller ./ opCaller

import (
	"context"
	"time"

	"github.com/opctl/sdk-golang/model"
	"github.com/opctl/sdk-golang/opspec/interpreter/call/op/outputs"
	"github.com/opctl/sdk-golang/opspec/opfile"
	"github.com/opctl/sdk-golang/util/pubsub"
)

type opCaller interface {
	// Executes an op call
	Call(
		dcgOpCall *model.DCGOpCall,
		inboundScope map[string]*model.Value,
		scgOpCall *model.SCGOpCall,
	) error
}

func newOpCaller(
	pubSub pubsub.PubSub,
	dcgNodeRepo dcgNodeRepo,
	caller caller,
	dataDirPath string,
) opCaller {
	return _opCaller{
		outputsInterpreter: outputs.NewInterpreter(),
		dotYmlGetter:       dotyml.NewGetter(),
		pubSub:             pubSub,
		dcgNodeRepo:        dcgNodeRepo,
		caller:             caller,
	}
}

type _opCaller struct {
	outputsInterpreter outputs.Interpreter
	dotYmlGetter       dotyml.Getter
	pubSub             pubsub.PubSub
	dcgNodeRepo        dcgNodeRepo
	caller             caller
}

func (oc _opCaller) Call(
	dcgOpCall *model.DCGOpCall,
	inboundScope map[string]*model.Value,
	scgOpCall *model.SCGOpCall,
) error {
	var err error
	var isKilled bool
	outputs := map[string]*model.Value{}

	defer func() {
		// defer must be defined before conditional return statements so it always runs
		if isKilled {
			// guard: op killed (we got preempted)
			oc.pubSub.Publish(
				model.Event{
					Timestamp: time.Now().UTC(),
					OpEnded: &model.OpEndedEvent{
						OpID:     dcgOpCall.OpID,
						Outcome:  model.OpOutcomeKilled,
						RootOpID: dcgOpCall.RootOpID,
						OpRef:    dcgOpCall.OpHandle.Ref(),
					},
				},
			)
			return
		}

		oc.dcgNodeRepo.DeleteIfExists(dcgOpCall.OpID)

		var opOutcome string
		if nil != err {
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
					Outputs:  outputs,
				},
			},
		)

	}()

	oc.dcgNodeRepo.Add(
		&dcgNodeDescriptor{
			Id:       dcgOpCall.OpID,
			OpRef:    dcgOpCall.OpHandle.Ref(),
			RootOpID: dcgOpCall.RootOpID,
			Op:       &dcgOpDescriptor{},
		},
	)

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
		dcgOpCall.RootOpID,
	)

	isKilled = nil == oc.dcgNodeRepo.GetIfExists(dcgOpCall.RootOpID)

	if nil != err {
		return err
	}

	// wait on op outputs
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

	// filter op outputs to bound call outputs
	for boundName, boundValue := range scgOpCall.Outputs {
		// return bound outputs
		if "" == boundValue {
			// implicit value
			boundValue = boundName
		}
		for opOutputName, opOutputValue := range opOutputs {
			if boundValue == opOutputName {
				outputs[boundName] = opOutputValue
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
			// parallel calls have no outputs
			return nil
		case nil != event.CallEnded && event.CallEnded.CallID == dcgOpCall.ChildCallID:
			break eventLoop
		}
	}

	return opOutputs
}

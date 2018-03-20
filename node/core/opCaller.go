package core

//go:generate counterfeiter -o ./fakeOpCaller.go --fake-name fakeOpCaller ./ opCaller

import (
	"context"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/op/dotyml"
	"github.com/opspec-io/sdk-golang/op/interpreter/opcall"
	"github.com/opspec-io/sdk-golang/op/interpreter/opcall/outputs"
	"github.com/opspec-io/sdk-golang/util/pubsub"
	"time"
)

type opCaller interface {
	// Executes an op call
	Call(
		inboundScope map[string]*model.Value,
		opId string,
		opDirHandle model.DataHandle,
		rootOpId string,
		scgOpCall *model.SCGOpCall,
	) error
}

func newOpCaller(
	pubSub pubsub.PubSub,
	dcgNodeRepo dcgNodeRepo,
	caller caller,
	rootFSPath string,
) opCaller {
	return _opCaller{
		outputsInterpreter: outputs.NewInterpreter(),
		opCallInterpreter:  opcall.NewInterpreter(rootFSPath),
		dotYmlGetter:       dotyml.NewGetter(),
		pubSub:             pubSub,
		dcgNodeRepo:        dcgNodeRepo,
		caller:             caller,
	}
}

type _opCaller struct {
	outputsInterpreter outputs.Interpreter
	opCallInterpreter  opcall.Interpreter
	dotYmlGetter       dotyml.Getter
	pubSub             pubsub.PubSub
	dcgNodeRepo        dcgNodeRepo
	caller             caller
}

func (oc _opCaller) Call(
	inboundScope map[string]*model.Value,
	opId string,
	opDirHandle model.DataHandle,
	rootOpId string,
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
						OpId:     opId,
						Outcome:  model.OpOutcomeKilled,
						RootOpId: rootOpId,
						PkgRef:   scgOpCall.Pkg.Ref,
					},
				},
			)
			return
		}

		oc.dcgNodeRepo.DeleteIfExists(opId)

		var opOutcome string
		if nil != err {
			oc.pubSub.Publish(
				model.Event{
					Timestamp: time.Now().UTC(),
					OpErred: &model.OpErredEvent{
						Msg:      err.Error(),
						OpId:     opId,
						PkgRef:   scgOpCall.Pkg.Ref,
						RootOpId: rootOpId,
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
					OpId:     opId,
					PkgRef:   scgOpCall.Pkg.Ref,
					Outcome:  opOutcome,
					RootOpId: rootOpId,
					Outputs:  outputs,
				},
			},
		)

	}()

	oc.dcgNodeRepo.Add(
		&dcgNodeDescriptor{
			Id:       opId,
			PkgRef:   scgOpCall.Pkg.Ref,
			RootOpId: rootOpId,
			Op:       &dcgOpDescriptor{},
		},
	)

	dcgOpCall, err := oc.opCallInterpreter.Interpret(
		inboundScope,
		scgOpCall,
		opId,
		opDirHandle,
		rootOpId,
	)
	if nil != err {
		return err
	}

	oc.pubSub.Publish(
		model.Event{
			Timestamp: time.Now().UTC(),
			OpStarted: &model.OpStartedEvent{
				OpId:     opId,
				PkgRef:   scgOpCall.Pkg.Ref,
				RootOpId: rootOpId,
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
		dcgOpCall.ChildCallId,
		dcgOpCall.Inputs,
		dcgOpCall.ChildCallSCG,
		dcgOpCall.DataHandle,
		rootOpId,
	)

	isKilled = nil == oc.dcgNodeRepo.GetIfExists(rootOpId)

	if nil != err {
		return err
	}

	// wait on op outputs
	opOutputs := <-opOutputsChan

	opDotYml, err := oc.dotYmlGetter.Get(
		context.TODO(),
		dcgOpCall.DataHandle,
	)
	if nil != err {
		return err
	}
	opPath := dcgOpCall.DataHandle.Path()
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
			Roots: []string{dcgOpCall.RootOpId},
			Since: &eventFilterSince,
		},
	)

	opOutputs := map[string]*model.Value{}
eventLoop:
	for event := range eventChannel {
		switch {
		case nil != event.OpEnded && event.OpEnded.OpId == dcgOpCall.OpId:
			// parent ended prematurely
			return nil
		case nil != event.OpEnded && event.OpEnded.OpId == dcgOpCall.ChildCallId:
			for name, value := range event.OpEnded.Outputs {
				opOutputs[name] = value
			}
			break eventLoop
		case nil != event.ContainerExited && event.ContainerExited.ContainerId == dcgOpCall.ChildCallId:
			for name, value := range event.ContainerExited.Outputs {
				opOutputs[name] = value
			}
			break eventLoop
		case nil != event.SerialCallEnded && event.SerialCallEnded.CallId == dcgOpCall.ChildCallId:
			for name, value := range event.SerialCallEnded.Outputs {
				opOutputs[name] = value
			}
			break eventLoop
		case nil != event.ParallelCallEnded && event.ParallelCallEnded.CallId == dcgOpCall.ChildCallId:
			// parallel calls have no outputs
			return nil
		}
	}

	return opOutputs
}

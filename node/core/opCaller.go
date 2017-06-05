package core

//go:generate counterfeiter -o ./fakeOpCaller.go --fake-name fakeOpCaller ./ opCaller

import (
	"github.com/opctl/opctl/util/pubsub"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/opcall"
	"time"
)

type opCaller interface {
	// Executes an op call
	Call(
		inboundScope map[string]*model.Value,
		opId string,
		pkgBasePath string,
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
		opCall:      opcall.New(rootFSPath),
		pubSub:      pubSub,
		dcgNodeRepo: dcgNodeRepo,
		caller:      caller,
	}
}

type _opCaller struct {
	opCall      opcall.OpCall
	pubSub      pubsub.PubSub
	dcgNodeRepo dcgNodeRepo
	caller      caller
}

func (oc _opCaller) Call(
	inboundScope map[string]*model.Value,
	opId string,
	pkgBasePath string,
	rootOpId string,
	scgOpCall *model.SCGOpCall,
) error {
	var err error
	var isKilled bool
	var outputs map[string]*model.Value
	if "" != scgOpCall.Ref {
		// fallback for deprecated pkg ref format
		scgOpCall.Pkg = &model.SCGOpCallPkg{
			Ref: scgOpCall.Ref,
		}
	}
	defer func() {
		// defer must be defined before conditional return statements so it always runs
		if isKilled {
			// guard: op killed (we got preempted)
			oc.pubSub.Publish(
				&model.Event{
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
				&model.Event{
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
			&model.Event{
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

	dcgOpCall, err := oc.opCall.Interpret(
		inboundScope,
		scgOpCall,
		opId,
		pkgBasePath,
		rootOpId,
	)
	if nil != err {
		return err
	}

	oc.pubSub.Publish(
		&model.Event{
			Timestamp: time.Now().UTC(),
			OpStarted: &model.OpStartedEvent{
				OpId:     opId,
				PkgRef:   scgOpCall.Pkg.Ref,
				RootOpId: rootOpId,
			},
		},
	)

	// interpret outputs
	outputsChan := make(chan map[string]*model.Value, 1)
	go func() {
		outputsChan <- oc.interpretOutputs(
			scgOpCall,
			dcgOpCall,
		)
	}()

	err = oc.caller.Call(
		dcgOpCall.ChildCallId,
		dcgOpCall.Inputs,
		dcgOpCall.ChildCallSCG,
		dcgOpCall.PkgRef,
		rootOpId,
	)

	isKilled = nil == oc.dcgNodeRepo.GetIfExists(rootOpId)

	if nil != err {
		return err
	}

	// wait on outputs
	outputs = <-outputsChan

	return nil
}

func (oc _opCaller) interpretOutputs(
	scgOpCall *model.SCGOpCall,
	dcgOpCall *model.DCGOpCall,
) map[string]*model.Value {

	// subscribe to events
	eventChannel := make(chan *model.Event, 150)
	eventFilterSince := time.Now().UTC()
	oc.pubSub.Subscribe(
		&model.EventFilter{
			RootOpIds: []string{dcgOpCall.RootOpId},
			Since:     &eventFilterSince,
		},
		eventChannel,
	)

	childOutputs := map[string]*model.Value{}
eventLoop:
	for event := range eventChannel {
		switch {
		case nil != event.OpEnded && event.OpEnded.OpId == dcgOpCall.OpId:
			// parent ended prematurely
			return nil
		case nil != event.OpEnded && event.OpEnded.OpId == dcgOpCall.ChildCallId:
			for name, value := range event.OpEnded.Outputs {
				childOutputs[name] = value
			}
			break eventLoop
		case nil != event.ContainerExited && event.ContainerExited.ContainerId == dcgOpCall.ChildCallId:
			for name, value := range event.ContainerExited.Outputs {
				childOutputs[name] = value
			}
			break eventLoop
		case nil != event.SerialCallEnded && event.SerialCallEnded.CallId == dcgOpCall.ChildCallId:
			for name, value := range event.SerialCallEnded.Outputs {
				childOutputs[name] = value
			}
			break eventLoop
		case nil != event.ParallelCallEnded && event.ParallelCallEnded.CallId == dcgOpCall.ChildCallId:
			// parallel calls have no outputs
			return nil
		}
	}

	outputs := map[string]*model.Value{}
	for boundName, boundValue := range scgOpCall.Outputs {
		// return bound outputs
		if "" == boundValue {
			// implicit value
			boundValue = boundName
		}
		for childOutputName, childOutputValue := range childOutputs {
			if boundValue == childOutputName {
				outputs[boundName] = childOutputValue
			}
		}
	}

	return outputs
}

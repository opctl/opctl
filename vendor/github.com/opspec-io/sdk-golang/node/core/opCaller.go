package core

//go:generate counterfeiter -o ./fakeOpCaller.go --fake-name fakeOpCaller ./ opCaller

import (
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/opcall"
	outputsPkg "github.com/opspec-io/sdk-golang/opcall/outputs"
	"github.com/opspec-io/sdk-golang/pkg"
	"github.com/opspec-io/sdk-golang/util/pubsub"
	"time"
)

type opCaller interface {
	// Executes an op call
	Call(
		inboundScope map[string]*model.Value,
		opId string,
		pkgHandle model.PkgHandle,
		rootOpId string,
		scgOpCall *model.SCGOpCall,
	) error
}

func newOpCaller(
	pubSub pubsub.PubSub,
	caller caller,
	rootFSPath string,
) opCaller {
	return _opCaller{
		outputs: outputsPkg.New(),
		opCall:  opcall.New(rootFSPath),
		pkg:     pkg.New(),
		pubSub:  pubSub,
		caller:  caller,
	}
}

type _opCaller struct {
	outputs outputsPkg.Outputs
	opCall  opcall.OpCall
	pkg     pkg.Pkg
	pubSub  pubsub.PubSub
	caller  caller
}

func (oc _opCaller) Call(
	inboundScope map[string]*model.Value,
	opId string,
	pkgHandle model.PkgHandle,
	rootOpId string,
	scgOpCall *model.SCGOpCall,
) error {
	var err error
	var isKilled bool
	var outputs map[string]*model.Value
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

	dcgOpCall, err := oc.opCall.Interpret(
		inboundScope,
		scgOpCall,
		opId,
		pkgHandle,
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
		dcgOpCall.PkgHandle,
		rootOpId,
	)

	isKilled = false

	if nil != err {
		return err
	}

	// wait on outputs
	outputs = <-outputsChan

	childPkg, err := oc.pkg.GetManifest(dcgOpCall.PkgHandle)
	if nil != err {
		return err
	}

	childPkgPath := dcgOpCall.PkgHandle.Path()
	outputs, err = oc.outputs.Interpret(outputs, childPkg.Outputs, *childPkgPath)

	return err
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
			Roots: []string{dcgOpCall.RootOpId},
			Since: &eventFilterSince,
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

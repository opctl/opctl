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
		inboundScope map[string]*model.Data,
		opId string,
		pkgBasePath string,
		rootOpId string,
		scgOpCall *model.SCGOpCall,
	) (
		err error,
	)
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

func (this _opCaller) Call(
	inboundScope map[string]*model.Data,
	opId string,
	pkgBasePath string,
	rootOpId string,
	scgOpCall *model.SCGOpCall,
) (
	err error,
) {
	if "" != scgOpCall.Ref {
		// fallback for deprecated pkg ref format
		scgOpCall.Pkg = &model.SCGOpCallPkg{
			Ref: scgOpCall.Ref,
		}
	}

	defer func() {
		// defer must be defined before conditional return statements so it always runs

		if nil == this.dcgNodeRepo.GetIfExists(rootOpId) {
			// guard: op killed (we got preempted)
			this.pubSub.Publish(
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

		this.dcgNodeRepo.DeleteIfExists(opId)

		var opOutcome string
		if nil != err {
			this.pubSub.Publish(
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

		this.pubSub.Publish(
			&model.Event{
				Timestamp: time.Now().UTC(),
				OpEnded: &model.OpEndedEvent{
					OpId:     opId,
					PkgRef:   scgOpCall.Pkg.Ref,
					Outcome:  opOutcome,
					RootOpId: rootOpId,
				},
			},
		)

	}()

	this.dcgNodeRepo.Add(
		&dcgNodeDescriptor{
			Id:       opId,
			PkgRef:   scgOpCall.Pkg.Ref,
			RootOpId: rootOpId,
			Op:       &dcgOpDescriptor{},
		},
	)

	dcgOpCall, err := this.opCall.Interpret(
		inboundScope,
		scgOpCall,
		opId,
		pkgBasePath,
		rootOpId,
	)
	if nil != err {
		return
	}

	this.pubSub.Publish(
		&model.Event{
			Timestamp: time.Now().UTC(),
			OpStarted: &model.OpStartedEvent{
				OpId:     opId,
				PkgRef:   scgOpCall.Pkg.Ref,
				RootOpId: rootOpId,
			},
		},
	)

	go this.txOutputs(dcgOpCall, scgOpCall)

	err = this.caller.Call(
		dcgOpCall.ChildCallId,
		dcgOpCall.ChildCallScope,
		dcgOpCall.ChildCallSCG,
		dcgOpCall.PkgRef,
		rootOpId,
	)
	if nil != err {
		return
	}

	return

}

func (this _opCaller) txOutputs(
	dcgOpCall *model.DCGOpCall,
	scgOpCall *model.SCGOpCall,
) {
	// subscribe to events
	eventChannel := make(chan *model.Event, 150)
	eventFilterSince := time.Now().UTC()
	this.pubSub.Subscribe(
		&model.EventFilter{
			RootOpIds: []string{dcgOpCall.RootOpId},
			Since:     &eventFilterSince,
		},
		eventChannel,
	)

	// send outputs
eventLoop:
	for event := range eventChannel {
		switch {
		case nil != event.OpEnded && event.OpEnded.OpId == dcgOpCall.OpId:
			break eventLoop
		case nil != event.OutputInitialized && event.OutputInitialized.CallId == dcgOpCall.ChildCallId:
			childOutput := event.OutputInitialized
			if _, ok := scgOpCall.Outputs[childOutput.Name]; ok {
				this.pubSub.Publish(&model.Event{
					Timestamp: time.Now().UTC(),
					OutputInitialized: &model.OutputInitializedEvent{
						Name:     childOutput.Name,
						Value:    childOutput.Value,
						RootOpId: dcgOpCall.RootOpId,
						CallId:   dcgOpCall.OpId,
					},
				})
			}
		}
	}
}

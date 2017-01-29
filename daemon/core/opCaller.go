package core

//go:generate counterfeiter -o ./fakeOpCaller.go --fake-name fakeOpCaller ./ opCaller

import (
	"github.com/opspec-io/opctl/util/eventbus"
	"github.com/opspec-io/opctl/util/uniquestring"
	"github.com/opspec-io/sdk-golang/pkg/bundle"
	"github.com/opspec-io/sdk-golang/pkg/model"
	"github.com/opspec-io/sdk-golang/pkg/validate"
	"strings"
	"time"
)

type opCaller interface {
	// Executes an op call
	Call(
		inboundScope map[string]*model.Data,
		opId string,
		opRef string,
		opGraphId string,
	) (
		outboundScope map[string]*model.Data,
		err error,
	)
}

func newOpCaller(
	bundle bundle.Bundle,
	eventBus eventbus.EventBus,
	dcgNodeRepo dcgNodeRepo,
	caller caller,
	uniqueStringFactory uniquestring.UniqueStringFactory,
	validate validate.Validate,
) opCaller {
	return _opCaller{
		bundle:              bundle,
		eventBus:            eventBus,
		dcgNodeRepo:         dcgNodeRepo,
		caller:              caller,
		uniqueStringFactory: uniqueStringFactory,
		validate:            validate,
	}
}

type _opCaller struct {
	bundle              bundle.Bundle
	eventBus            eventbus.EventBus
	dcgNodeRepo         dcgNodeRepo
	caller              caller
	uniqueStringFactory uniquestring.UniqueStringFactory
	validate            validate.Validate
}

func (this _opCaller) Call(
	inboundScope map[string]*model.Data,
	opId string,
	opRef string,
	opGraphId string,
) (
	outboundScope map[string]*model.Data,
	err error,
) {

	this.dcgNodeRepo.Add(
		&dcgNodeDescriptor{
			Id:        opId,
			OpRef:     opRef,
			OpGraphId: opGraphId,
			Op:        &dcgOpDescriptor{},
		},
	)

	op, err := this.bundle.GetOp(
		opRef,
	)
	if nil != err {
		this.eventBus.Publish(
			model.Event{
				Timestamp: time.Now().UTC(),
				OpEncounteredError: &model.OpEncounteredErrorEvent{
					Msg:       err.Error(),
					OpId:      opId,
					OpRef:     opRef,
					OpGraphId: opGraphId,
				},
			},
		)
		return
	}

	// validate args
	errs := []error{}
	for _, input := range op.Inputs {
		var arg *model.Data

		// resolve var for input
		switch {
		case nil != input.Dir:
			arg = inboundScope[input.Dir.Name]
		case nil != input.File:
			arg = inboundScope[input.File.Name]
		case nil != input.Socket:
			arg = inboundScope[input.Socket.Name]
		case nil != input.String:
			arg = inboundScope[input.String.Name]
		}
		errs = append(errs, this.validate.Param(arg, input)...)
	}
	if len(errs) > 0 {
		errStrings := []string{}
		for _, err := range errs {
			errStrings = append(errStrings, err.Error())
		}
		this.eventBus.Publish(
			model.Event{
				Timestamp: time.Now().UTC(),
				OpEncounteredError: &model.OpEncounteredErrorEvent{
					Msg:       strings.Join(errStrings, "\n"),
					OpId:      opId,
					OpRef:     opRef,
					OpGraphId: opGraphId,
				},
			},
		)
		return
	}

	opStartedEvent := model.Event{
		Timestamp: time.Now().UTC(),
		OpStarted: &model.OpStartedEvent{
			OpId:      opId,
			OpRef:     opRef,
			OpGraphId: opGraphId,
		},
	}
	this.eventBus.Publish(opStartedEvent)

	outboundScope, err = this.caller.Call(
		this.uniqueStringFactory.Construct(),
		inboundScope,
		op.Run,
		opRef,
		opGraphId,
	)

	defer func(err error) {

		if nil == this.dcgNodeRepo.GetIfExists(opGraphId) {
			// guard: op killed (we got preempted)
			this.eventBus.Publish(
				model.Event{
					Timestamp: time.Now().UTC(),
					OpEnded: &model.OpEndedEvent{
						OpId:      opId,
						Outcome:   model.OpOutcomeKilled,
						OpGraphId: opGraphId,
						OpRef:     opRef,
					},
				},
			)
			return
		}

		this.dcgNodeRepo.DeleteIfExists(opId)

		var opOutcome string
		if nil != err {
			opOutcome = model.OpOutcomeFailed
			this.eventBus.Publish(
				model.Event{
					Timestamp: time.Now().UTC(),
					OpEncounteredError: &model.OpEncounteredErrorEvent{
						Msg:       err.Error(),
						OpId:      opId,
						OpRef:     opRef,
						OpGraphId: opGraphId,
					},
				},
			)
		} else {
			opOutcome = model.OpOutcomeSucceeded
		}

		this.eventBus.Publish(
			model.Event{
				Timestamp: time.Now().UTC(),
				OpEnded: &model.OpEndedEvent{
					OpId:      opId,
					OpRef:     opRef,
					Outcome:   opOutcome,
					OpGraphId: opGraphId,
				},
			},
		)

	}(err)

	return

}

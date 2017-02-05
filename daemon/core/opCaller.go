package core

//go:generate counterfeiter -o ./fakeOpCaller.go --fake-name fakeOpCaller ./ opCaller

import (
	"github.com/opspec-io/opctl/util/pubsub"
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
	pubSub pubsub.PubSub,
	dcgNodeRepo dcgNodeRepo,
	caller caller,
	uniqueStringFactory uniquestring.UniqueStringFactory,
	validate validate.Validate,
) opCaller {
	return _opCaller{
		bundle:              bundle,
		pubSub:              pubSub,
		dcgNodeRepo:         dcgNodeRepo,
		caller:              caller,
		uniqueStringFactory: uniqueStringFactory,
		validate:            validate,
	}
}

type _opCaller struct {
	bundle              bundle.Bundle
	pubSub              pubsub.PubSub
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
		this.pubSub.Publish(
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
	for inputName, input := range op.Inputs {
		var arg *model.Data

		// resolve var for input
		switch {
		case nil != input.Dir:
			arg = inboundScope[inputName]
		case nil != input.File:
			arg = inboundScope[inputName]
		case nil != input.Socket:
			arg = inboundScope[inputName]
		case nil != input.String:
			if inScopeVar, ok := inboundScope[inputName]; ok {
				arg = inScopeVar
			} else {
				// fallback to default
				arg = &model.Data{String: input.String.Default}
			}
		}
		errs = append(errs, this.validate.Param(arg, input)...)
	}
	if len(errs) > 0 {
		errStrings := []string{}
		for _, err := range errs {
			errStrings = append(errStrings, err.Error())
		}
		this.pubSub.Publish(
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
	this.pubSub.Publish(opStartedEvent)

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
			this.pubSub.Publish(
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
			this.pubSub.Publish(
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

		this.pubSub.Publish(
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

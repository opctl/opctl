package core

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
	Call(
		args map[string]*model.Data,
		opId string,
		opRef string,
		opGraphId string,
	) (
		outputs map[string]*model.Data,
		err error,
	)
}

type _opCaller struct {
	bundle              bundle.Bundle
	eventBus            eventbus.EventBus
	nodeRepo            nodeRepo
	caller              caller
	uniqueStringFactory uniquestring.UniqueStringFactory
	validate            validate.Validate
}

// Calls an op
func (this _opCaller) Call(
	args map[string]*model.Data,
	opId string,
	opRef string,
	opGraphId string,
) (
	outputs map[string]*model.Data,
	err error,
) {

	this.nodeRepo.add(
		&nodeDescriptor{
			Id:        opId,
			OpRef:     opRef,
			OpGraphId: opGraphId,
			Op:        &opDescriptor{},
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

		switch {
		case nil != input.String:
			if providedArg := args[input.String.Name]; nil != providedArg {
				arg = providedArg
			}
		case nil != input.NetSocket:
			if providedArg := args[input.NetSocket.Name]; nil != providedArg {
				arg = providedArg
			}
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

	outputs, err = this.caller.Call(
		this.uniqueStringFactory.Construct(),
		args,
		op.Run,
		opRef,
		opGraphId,
	)

	defer func(err error) {

		if nil == this.nodeRepo.getIfExists(opGraphId) {
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

		this.nodeRepo.deleteIfExists(opId)

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

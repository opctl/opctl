package core

//go:generate counterfeiter -o ./fakeOpCaller.go --fake-name fakeOpCaller ./ opCaller

import (
	"bytes"
	"fmt"
	"github.com/opspec-io/opctl/util/pubsub"
	"github.com/opspec-io/opctl/util/uniquestring"
	"github.com/opspec-io/sdk-golang/pkg/managepackages"
	"github.com/opspec-io/sdk-golang/pkg/model"
	"github.com/opspec-io/sdk-golang/pkg/validate"
	"strconv"
	"time"
)

type opCaller interface {
	// Executes an op call
	Call(
		inputScope map[string]*model.Data,
		outputs chan *variable,
		opId string,
		pkgRef string,
		rootOpId string,
	) (
		err error,
	)
}

func newOpCaller(
	pkg managepackages.ManagePackages,
	pubSub pubsub.PubSub,
	dcgNodeRepo dcgNodeRepo,
	caller caller,
	uniqueStringFactory uniquestring.UniqueStringFactory,
	validate validate.Validate,
) opCaller {
	return _opCaller{
		managePackages:      pkg,
		pubSub:              pubSub,
		dcgNodeRepo:         dcgNodeRepo,
		caller:              caller,
		uniqueStringFactory: uniqueStringFactory,
		validate:            validate,
	}
}

type _opCaller struct {
	managePackages      managepackages.ManagePackages
	pubSub              pubsub.PubSub
	dcgNodeRepo         dcgNodeRepo
	caller              caller
	uniqueStringFactory uniquestring.UniqueStringFactory
	validate            validate.Validate
}

func (this _opCaller) Call(
	inputScope map[string]*model.Data,
	outputs chan *variable,
	opId string,
	pkgRef string,
	rootOpId string,
) (
	err error,
) {
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
						PkgRef:   pkgRef,
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
					OpEncounteredError: &model.OpEncounteredErrorEvent{
						Msg:      err.Error(),
						OpId:     opId,
						PkgRef:   pkgRef,
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
					PkgRef:   pkgRef,
					Outcome:  opOutcome,
					RootOpId: rootOpId,
				},
			},
		)

	}()

	this.dcgNodeRepo.Add(
		&dcgNodeDescriptor{
			Id:       opId,
			PkgRef:   pkgRef,
			RootOpId: rootOpId,
			Op:       &dcgOpDescriptor{},
		},
	)

	op, err := this.managePackages.GetPackage(
		pkgRef,
	)
	if nil != err {
		return
	}

	this.applyParamDefaultsToScope(inputScope, op.Inputs)

	// validate inputs
	err = this.validateScope("input", inputScope, op.Inputs)
	if nil != err {
		return
	}

	this.pubSub.Publish(
		&model.Event{
			Timestamp: time.Now().UTC(),
			OpStarted: &model.OpStartedEvent{
				OpId:     opId,
				PkgRef:   pkgRef,
				RootOpId: rootOpId,
			},
		},
	)

	childOutputs := make(chan *variable, 150)
	outputScope := map[string]*model.Data{}
	// stream outputs as they're sent
	go func() {
		for childOutput := range childOutputs {
			varName := childOutput.Name
			varValue := childOutput.Value
			if param, ok := op.Outputs[varName]; ok {

				childOutput = &variable{
					Name:  varName,
					Value: this.getValueOrDefault(varName, varValue, param),
				}

				// track outputs in outputScope so we can validate them before returning
				outputScope[childOutput.Name] = childOutput.Value
				outputs <- childOutput
			}
		}
		close(outputs)
	}()

	err = this.caller.Call(
		this.uniqueStringFactory.Construct(),
		inputScope,
		childOutputs,
		op.Run,
		pkgRef,
		rootOpId,
	)
	if nil != err {
		return
	}

	// validate outputs
	err = this.validateScope("output", outputScope, op.Outputs)

	return

}

func (this _opCaller) applyParamDefaultsToScope(
	scope map[string]*model.Data,
	params map[string]*model.Param,
) {
	for paramName, param := range params {
		this.getValueOrDefault(paramName, scope[paramName], param)
	}
}

func (this _opCaller) getValueOrDefault(
	varName string,
	varValue *model.Data,
	param *model.Param,
) *model.Data {
	switch {
	case nil != param.Number:
		if nil == varValue {
			return &model.Data{Number: param.Number.Default}
		}
	case nil != param.String:
		if nil == varValue {
			// apply default; value not found in scope
			return &model.Data{String: param.String.Default}
		}
	}
	return varValue
}

func (this _opCaller) validateParam(
	paramType string,
	varName string,
	varValue *model.Data,
	param *model.Param,
) error {
	var (
		argDisplayValue string
	)

	if nil != varValue {
		switch {
		case nil != param.Dir:
			argDisplayValue = varValue.Dir
		case nil != param.File:
			argDisplayValue = varValue.File
		case nil != param.Number:
			if param.Number.IsSecret {
				argDisplayValue = "************"
			} else {
				argDisplayValue = strconv.FormatFloat(varValue.Number, 'f', -1, 64)
			}
		case nil != param.Socket:
			argDisplayValue = varValue.Socket
		case nil != param.String:
			if param.String.IsSecret {
				argDisplayValue = "************"
			} else {
				argDisplayValue = varValue.String
			}
		}
	}

	// validate
	validationErrors := this.validate.Param(varValue, param)

	messageBuffer := bytes.NewBufferString(``)

	if len(validationErrors) > 0 {
		messageBuffer.WriteString(fmt.Sprintf(`
  Name: %v
  Value: %v
  Error(s):`, varName, argDisplayValue),
		)
		for _, validationError := range validationErrors {
			messageBuffer.WriteString(fmt.Sprintf(`
    - %v`, validationError.Error()))
		}
		messageBuffer.WriteString(`
`)
	}

	if messageBuffer.Len() > 0 {
		return fmt.Errorf(`
-
  validation of the following %v failed:
%v
-`, paramType, messageBuffer.String())
	}
	return nil
}

func (this _opCaller) validateScope(
	scopeType string,
	scope map[string]*model.Data,
	params map[string]*model.Param,
) error {
	for paramName, param := range params {
		if err := this.validateParam(scopeType, paramName, scope[paramName], param); nil != err {
			return err
		}
	}
	return nil
}

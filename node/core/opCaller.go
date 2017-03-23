package core

//go:generate counterfeiter -o ./fakeOpCaller.go --fake-name fakeOpCaller ./ opCaller

import (
	"bytes"
	"fmt"
	"github.com/opctl/opctl/util/pubsub"
	"github.com/opctl/opctl/util/uniquestring"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/pkg"
	"github.com/opspec-io/sdk-golang/validate"
	"github.com/pkg/errors"
	"strconv"
	"time"
)

type opCaller interface {
	// Executes an op call
	Call(
		inboundScope map[string]*model.Data,
		opId string,
		pkgRef string,
		rootOpId string,
		scgOpCall *model.ScgOpCall,
	) (
		err error,
	)
}

func newOpCaller(
	pkg pkg.Pkg,
	pubSub pubsub.PubSub,
	dcgNodeRepo dcgNodeRepo,
	caller caller,
	uniqueStringFactory uniquestring.UniqueStringFactory,
	validate validate.Validate,
) opCaller {
	return _opCaller{
		pkg:                 pkg,
		pubSub:              pubSub,
		dcgNodeRepo:         dcgNodeRepo,
		caller:              caller,
		uniqueStringFactory: uniqueStringFactory,
		validate:            validate,
	}
}

type _opCaller struct {
	pkg                 pkg.Pkg
	pubSub              pubsub.PubSub
	dcgNodeRepo         dcgNodeRepo
	caller              caller
	uniqueStringFactory uniquestring.UniqueStringFactory
	validate            validate.Validate
}

func (this _opCaller) Call(
	inboundScope map[string]*model.Data,
	opId string,
	pkgRef string,
	rootOpId string,
	scgOpCall *model.ScgOpCall,
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

	dcgOpCall := &model.DCGOpCall{
		DCGBaseCall: &model.DCGBaseCall{
			RootOpId: rootOpId,
			PkgRef:   pkgRef,
		},
		OpId:        opId,
		ChildCallId: this.uniqueStringFactory.Construct(),
	}

	pkg, err := this.pkg.Get(
		pkgRef,
	)
	if nil != err {
		return
	}

	inputs := map[string]*model.Data{}
	for inputName, scopeName := range scgOpCall.Inputs {
		if "" == scopeName {
			// when not explicitly provided, set scopeName to inputName
			scopeName = inputName
		}
		inputs[inputName] = inboundScope[scopeName]
	}

	this.applyParamDefaultsToScope(inputs, pkg.Inputs)

	// validate inputs
	err = this.validateScope("input", inputs, pkg.Inputs)
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

	go this.txOutputs(dcgOpCall, &pkg)

	callErr := this.caller.Call(
		dcgOpCall.ChildCallId,
		inputs,
		pkg.Run,
		pkgRef,
		rootOpId,
	)
	if nil != callErr {
		err = errors.New("Error encountered during call")
		return
	}

	return

}

func (this _opCaller) txOutputs(
	dcgOpCall *model.DCGOpCall,
	pkg *model.PackageView,
) {
	// subscribe to events
	eventChannel := make(chan *model.Event, 150)
	this.pubSub.Subscribe(
		&model.EventFilter{RootOpIds: []string{dcgOpCall.RootOpId}},
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
			if _, ok := pkg.Outputs[childOutput.Name]; ok {
				this.pubSub.Publish(&model.Event{
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

func (this _opCaller) applyParamDefaultsToScope(
	scope map[string]*model.Data,
	params map[string]*model.Param,
) {
	for paramName, param := range params {
		scope[paramName] = this.getValueOrDefault(paramName, scope[paramName], param)
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
	messageBuffer := bytes.NewBufferString("")
	for paramName, param := range params {
		if err := this.validateParam(scopeType, paramName, scope[paramName], param); nil != err {
			messageBuffer.WriteString(err.Error())
		}
	}
	if messageBuffer.Len() > 0 {
		return errors.New(messageBuffer.String())
	}
	return nil
}

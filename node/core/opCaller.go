package core

//go:generate counterfeiter -o ./fakeOpCaller.go --fake-name fakeOpCaller ./ opCaller

import (
	"bytes"
	"fmt"
	"github.com/opctl/opctl/util/pubsub"
	"github.com/opctl/opctl/util/uniquestring"
	interpolatePkg "github.com/opspec-io/sdk-golang/interpolater"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/pkg"
	"github.com/opspec-io/sdk-golang/pkg/manifest"
	"github.com/opspec-io/sdk-golang/validate"
	"github.com/pkg/errors"
	"path/filepath"
	"strconv"
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
	manifest manifest.Manifest,
	pkg pkg.Pkg,
	pubSub pubsub.PubSub,
	dcgNodeRepo dcgNodeRepo,
	caller caller,
	uniqueStringFactory uniquestring.UniqueStringFactory,
	validate validate.Validate,
	rootFSPath string,
) opCaller {
	return _opCaller{
		manifest:            manifest,
		pkg:                 pkg,
		pubSub:              pubSub,
		dcgNodeRepo:         dcgNodeRepo,
		caller:              caller,
		uniqueStringFactory: uniqueStringFactory,
		validate:            validate,
		pkgCachePath:        filepath.Join(rootFSPath, "pkgs"),
	}
}

type _opCaller struct {
	manifest            manifest.Manifest
	pkg                 pkg.Pkg
	pubSub              pubsub.PubSub
	dcgNodeRepo         dcgNodeRepo
	caller              caller
	uniqueStringFactory uniquestring.UniqueStringFactory
	validate            validate.Validate
	pkgCachePath        string
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

	var pkgPath string
	if "" != scgOpCall.Ref {
		// handle deprecated pkgRef format
		pkgPath = filepath.Join(pkgBasePath, scgOpCall.Ref)
	} else if pkgPath, err = this.getPkgPath(
		pkgBasePath,
		inboundScope,
		scgOpCall.Pkg,
	); nil != err {
		return
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
						PkgRef:   pkgPath,
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
						PkgRef:   pkgPath,
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
					PkgRef:   pkgPath,
					Outcome:  opOutcome,
					RootOpId: rootOpId,
				},
			},
		)

	}()

	this.dcgNodeRepo.Add(
		&dcgNodeDescriptor{
			Id:       opId,
			PkgRef:   pkgPath,
			RootOpId: rootOpId,
			Op:       &dcgOpDescriptor{},
		},
	)

	dcgOpCall := &model.DCGOpCall{
		DCGBaseCall: &model.DCGBaseCall{
			RootOpId: rootOpId,
			PkgRef:   pkgPath,
		},
		OpId:        opId,
		ChildCallId: this.uniqueStringFactory.Construct(),
	}

	_pkg, err := this.manifest.Unmarshal(filepath.Join(pkgPath, pkg.OpDotYmlFileName))
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

	this.applyParamDefaultsToScope(inputs, _pkg.Inputs)

	// validate inputs
	err = this.validateScope("input", inputs, _pkg.Inputs)
	if nil != err {
		return
	}

	this.pubSub.Publish(
		&model.Event{
			Timestamp: time.Now().UTC(),
			OpStarted: &model.OpStartedEvent{
				OpId:     opId,
				PkgRef:   pkgPath,
				RootOpId: rootOpId,
			},
		},
	)

	go this.txOutputs(dcgOpCall, scgOpCall)

	err = this.caller.Call(
		dcgOpCall.ChildCallId,
		inputs,
		_pkg.Run,
		pkgPath,
		rootOpId,
	)
	if nil != err {
		return
	}

	return

}

func (this _opCaller) getPkgPath(
	pkgBasePath string,
	inboundScope map[string]*model.Data,
	scgOpCallPkg *model.SCGOpCallPkg,
) (string, error) {
	pkgRef, err := this.pkg.ParseRef(scgOpCallPkg.Ref)
	if nil != err {
		return "", err
	}

	interpolater := interpolatePkg.New()
	var username, password string
	if scgPullCreds := scgOpCallPkg.PullCreds; nil != scgPullCreds {
		username = interpolater.Interpolate(scgPullCreds.Username, inboundScope)
		password = interpolater.Interpolate(scgPullCreds.Password, inboundScope)
	}

	pkgPath, ok := this.pkg.Resolve(pkgRef, pkgBasePath, this.pkgCachePath)
	if !ok {
		// pkg not resolved; attempt to pull it
		err := this.pkg.Pull(this.pkgCachePath, pkgRef, &pkg.PullOpts{Username: username, Password: password})
		if nil != err {
			return "", err
		}

		// resolve pulled pkg
		pkgPath, ok = this.pkg.Resolve(pkgRef, pkgBasePath, this.pkgCachePath)
		if !ok {
			return "", fmt.Errorf("Unable to resolve pulled pkg '%v' from '%v'", pkgRef, this.pkgCachePath)
		}
	}
	return pkgPath, nil
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
			argDisplayValue = *varValue.Dir
		case nil != param.File:
			argDisplayValue = *varValue.File
		case nil != param.Number:
			if param.Number.IsSecret {
				argDisplayValue = "************"
			} else {
				argDisplayValue = strconv.FormatFloat(*varValue.Number, 'f', -1, 64)
			}
		case nil != param.Socket:
			argDisplayValue = *varValue.Socket
		case nil != param.String:
			if param.String.IsSecret {
				argDisplayValue = "************"
			} else {
				argDisplayValue = *varValue.String
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

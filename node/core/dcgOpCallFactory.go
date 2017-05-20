package core

//go:generate counterfeiter -o ./fakeDCGOpCallFactory.go --fake-name fakeDCGOpCallFactory ./ dcgOpCallFactory

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/opctl/opctl/util/uniquestring"
	"github.com/opspec-io/sdk-golang/interpolater"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/pkg"
	"github.com/opspec-io/sdk-golang/pkg/manifest"
	"github.com/opspec-io/sdk-golang/validate"
	"path/filepath"
	"strconv"
)

type dcgOpCallFactory interface {
	// Construct constructs a DCGOpCall
	Construct(
		scope map[string]*model.Data,
		scgOpCall *model.SCGOpCall,
		opId string,
		pkgBasePath string,
		rootOpId string,
	) (*model.DCGOpCall, error)
}

func newDCGOpCallFactory(
	rootFSPath string,
) dcgOpCallFactory {
	return _dcgOpCallFactory{
		interpolater:        interpolater.New(),
		manifest:            manifest.New(),
		pkg:                 pkg.New(),
		pkgCachePath:        filepath.Join(rootFSPath, "pkgs"),
		uniqueStringFactory: uniquestring.NewUniqueStringFactory(),
		validate:            validate.New(),
	}
}

type _dcgOpCallFactory struct {
	interpolater        interpolater.Interpolater
	manifest            manifest.Manifest
	pkg                 pkg.Pkg
	pkgCachePath        string
	uniqueStringFactory uniquestring.UniqueStringFactory
	validate            validate.Validate
}

func (this _dcgOpCallFactory) Construct(
	scope map[string]*model.Data,
	scgOpCall *model.SCGOpCall,
	opId string,
	pkgBasePath string,
	rootOpId string,
) (*model.DCGOpCall, error) {

	pkgPath, err := this.getPkgPath(
		pkgBasePath,
		scope,
		scgOpCall.Pkg,
	)
	if nil != err {
		return nil, err
	}

	_pkg, err := this.manifest.Unmarshal(filepath.Join(pkgPath, pkg.OpDotYmlFileName))
	if nil != err {
		return nil, err
	}

	dcgOpCall := &model.DCGOpCall{
		DCGBaseCall: &model.DCGBaseCall{
			RootOpId: rootOpId,
			PkgRef:   pkgPath,
		},
		OpId:           opId,
		ChildCallId:    this.uniqueStringFactory.Construct(),
		ChildCallSCG:   _pkg.Run,
		ChildCallScope: map[string]*model.Data{},
	}

	for inputName, scopeRef := range scgOpCall.Inputs {
		if "" == scopeRef {
			// when not explicitly provided, use inputName as scopeRef
			scopeRef = inputName
		}
		if scopeVal, ok := scope[scopeRef]; ok {
			dcgOpCall.ChildCallScope[inputName] = scopeVal
		}
	}

	this.applyParamDefaultsToScope(dcgOpCall.ChildCallScope, _pkg.Inputs)

	// validate inputs
	err = this.validateScope("input", dcgOpCall.ChildCallScope, _pkg.Inputs)
	if nil != err {
		return nil, err
	}

	return dcgOpCall, nil
}

func (this _dcgOpCallFactory) getPkgPath(
	pkgBasePath string,
	scope map[string]*model.Data,
	scgOpCallPkg *model.SCGOpCallPkg,
) (string, error) {
	pkgRef, err := this.pkg.ParseRef(scgOpCallPkg.Ref)
	if nil != err {
		return "", err
	}

	var username, password string
	if scgPullCreds := scgOpCallPkg.PullCreds; nil != scgPullCreds {
		username = this.interpolater.Interpolate(scgPullCreds.Username, scope)
		password = this.interpolater.Interpolate(scgPullCreds.Password, scope)
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

func (this _dcgOpCallFactory) applyParamDefaultsToScope(
	scope map[string]*model.Data,
	params map[string]*model.Param,
) {
	for paramName, param := range params {
		scope[paramName] = this.getValueOrDefault(paramName, scope[paramName], param)
	}
}

func (this _dcgOpCallFactory) getValueOrDefault(
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

func (this _dcgOpCallFactory) validateParam(
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

func (this _dcgOpCallFactory) validateScope(
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

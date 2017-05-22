package opcall

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/pkg"
	"path/filepath"
	"strconv"
	"strings"
)

func (this _OpCall) Interpret(
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
		ChildCallId:    this.uuid.NewV4().String(),
		ChildCallSCG:   _pkg.Run,
		ChildCallScope: map[string]*model.Data{},
	}

	for inputName, scgInputVal := range scgOpCall.Inputs {
		fmt.Printf("inputName: '%v'\nscgInputVal: '%v'\n", inputName, scgInputVal)
		switch {
		case "" == scgInputVal:
			fmt.Printf("empty: '%v'\n", scgInputVal)
			// bind implicit scopeRef
			if scopeVal, ok := scope[inputName]; ok {
				dcgOpCall.ChildCallScope[inputName] = scopeVal
			}
		case "" != this.scopeRef(scgInputVal):
			fmt.Printf("scopeRef: '%v'\n", scgInputVal)
			// bind explicit scopeRef
			if scopeVal, ok := scope[this.scopeRef(scgInputVal)]; ok {
				dcgOpCall.ChildCallScope[inputName] = scopeVal
			}
		default:
			fmt.Printf("default: '%v'\n", scgInputVal)
			// interpolate
			interpolatedVal := this.interpolater.Interpolate(scgInputVal, scope)
			if floatVal, err := strconv.ParseFloat(interpolatedVal, 64); nil == err {
				// bind number
				dcgOpCall.ChildCallScope[inputName] = &model.Data{Number: &floatVal}
				continue
			}
			dcgOpCall.ChildCallScope[inputName] = &model.Data{String: &interpolatedVal}
		}
	}

	this.applyParamDefaultsToScope(dcgOpCall.ChildCallScope, _pkg.Inputs)

	// paramvalidator inputs
	err = this.validateScope("input", dcgOpCall.ChildCallScope, _pkg.Inputs)
	if nil != err {
		return nil, err
	}

	return dcgOpCall, nil
}

// scopeRef gets a scopeRef from s; if unsuccessful returns empty string
func (this _OpCall) scopeRef(s string) string {
	i := strings.Index(s, "$(")
	if i >= 0 {
		j := strings.Index(s[i:], ")")
		if j >= 0 {
			return s[i+1 : j-i]
		}
	}
	return ""
}

func (this _OpCall) getPkgPath(
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

func (this _OpCall) applyParamDefaultsToScope(
	scope map[string]*model.Data,
	params map[string]*model.Param,
) {
	for paramName, param := range params {
		scope[paramName] = this.getValueOrDefault(paramName, scope[paramName], param)
	}
}

func (this _OpCall) getValueOrDefault(
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

func (this _OpCall) validateParam(
	paramType string,
	varName string,
	varValue *model.Data,
	param *model.Param,
) error {

	fmt.Printf("paramType:'%#v'\nvarName:'%#v'\nvarValue:'%#v'\nparam'%#v'\n", paramType, varName, varValue, param)
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

	// paramvalidator
	validationErrors := this.validate.Validate(varValue, param)

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

func (this _OpCall) validateScope(
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

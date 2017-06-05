package inputs

import (
	"github.com/opspec-io/sdk-golang/model"
	"strings"
)

func (_inputs _Inputs) Interpret(
	inputArgs map[string]string,
	inputParams map[string]*model.Param,
	pkgPath string,
	scope map[string]*model.Value,
) (map[string]*model.Value, []error) {
	dcgOpCallInputs := map[string]*model.Value{}

	for paramName, paramValue := range inputParams {
		// apply defaults
		if nil != paramValue {
			switch {
			case nil != paramValue.String && nil != paramValue.String.Default:
				dcgOpCallInputs[paramName] = &model.Value{String: paramValue.String.Default}
			case nil != paramValue.Number && nil != paramValue.Number.Default:
				dcgOpCallInputs[paramName] = &model.Value{Number: paramValue.Number.Default}
			case nil != paramValue.Dir && nil != paramValue.Dir.Default && strings.HasPrefix(*paramValue.Dir.Default, "/"):
				dcgOpCallInputs[paramName] = &model.Value{Dir: paramValue.Dir.Default}
			case nil != paramValue.File && nil != paramValue.File.Default && strings.HasPrefix(*paramValue.File.Default, "/"):
				dcgOpCallInputs[paramName] = &model.Value{File: paramValue.File.Default}
			}
		}
	}

	var errs []error
	for argName, argValue := range inputArgs {
		// override defaults w/ args
		var err error
		dcgOpCallInputs[argName], err = _inputs.argInterpreter.Interpret(
			argName,
			argValue,
			inputParams[argName],
			pkgPath,
			scope,
		)
		if nil != err {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return nil, errs
	}

	for inputName, inputValue := range dcgOpCallInputs {
		// validate inputs
		errs = append(errs, _inputs.validator.Validate(inputValue, inputParams[inputName])...)
	}

	return dcgOpCallInputs, errs
}

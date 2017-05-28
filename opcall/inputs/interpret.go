package inputs

import (
	"github.com/opspec-io/sdk-golang/model"
)

func (_inputs _Inputs) Interpret(
	inputArgs map[string]string,
	inputParams map[string]*model.Param,
	scope map[string]*model.Data,
) (map[string]*model.Data, []error) {
	dcgOpCallInputs := map[string]*model.Data{}

	for paramName, paramValue := range inputParams {
		// apply defaults
		if nil != paramValue {
			switch {
			case nil != paramValue.String && nil != paramValue.String.Default:
				dcgOpCallInputs[paramName] = &model.Data{String: paramValue.String.Default}
			case nil != paramValue.Number && nil != paramValue.Number.Default:
				dcgOpCallInputs[paramName] = &model.Data{Number: paramValue.Number.Default}
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
		_inputs.validator.Validate(inputValue, inputParams[inputName])
	}

	return dcgOpCallInputs, errs
}

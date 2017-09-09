package inputs

import (
	"github.com/opspec-io/sdk-golang/model"
	"path/filepath"
	"strings"
)

func (_inputs _Inputs) Interpret(
	inputArgs map[string]interface{},
	inputParams map[string]*model.Param,
	parentPkgRef,
	pkgRef string,
	scope map[string]*model.Value,
) (map[string]*model.Value, []error) {
	dcgOpCallInputs := map[string]*model.Value{}

	for paramName, paramValue := range inputParams {
		// apply defaults
		if nil != paramValue {
			switch {
			case nil != paramValue.Dir && nil != paramValue.Dir.Default && strings.HasPrefix(*paramValue.Dir.Default, "/"):
				dirValue := filepath.Join(pkgRef, *paramValue.Dir.Default)
				dcgOpCallInputs[paramName] = &model.Value{Dir: &dirValue}
			case nil != paramValue.File && nil != paramValue.File.Default && strings.HasPrefix(*paramValue.File.Default, "/"):
				fileValue := filepath.Join(pkgRef, *paramValue.File.Default)
				dcgOpCallInputs[paramName] = &model.Value{File: &fileValue}
			case nil != paramValue.Number && nil != paramValue.Number.Default:
				dcgOpCallInputs[paramName] = &model.Value{Number: paramValue.Number.Default}
			case nil != paramValue.Object && nil != paramValue.Object.Default:
				dcgOpCallInputs[paramName] = &model.Value{Object: paramValue.Object.Default}
			case nil != paramValue.String && nil != paramValue.String.Default:
				dcgOpCallInputs[paramName] = &model.Value{String: paramValue.String.Default}
			}
		}
	}

	var errs []error
	for argName, argValue := range inputArgs {
		// override defaults w/ args
		var err error
		// @TODO: argInterpreter.Interpret should perform validation so we don't need to below
		dcgOpCallInputs[argName], err = _inputs.argInterpreter.Interpret(
			argName,
			argValue,
			inputParams[argName],
			parentPkgRef,
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
		errs = append(errs, _inputs.data.Validate(inputValue, inputParams[inputName])...)
	}

	return dcgOpCallInputs, errs
}

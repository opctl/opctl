package inputs

import (
	"github.com/opspec-io/sdk-golang/data"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/opcall/params"
)

type interpreter interface {
	// Interpret interprets inputs via the provided inputArgs, inputParams, and scope;
	// opScratchDir will be used to store any run data such as scope coercions to files
	Interpret(
		inputArgs map[string]interface{},
		inputParams map[string]*model.Param,
		parentPkgHandle model.PkgHandle,
		pkgPath string,
		scope map[string]*model.Value,
		opScratchDir string,
	) (map[string]*model.Value, []error)
}

func newInterpreter() interpreter {
	return _interpreter{
		argInterpreter: newArgInterpreter(),
		data:           data.New(),
		params:         params.New(),
	}
}

type _interpreter struct {
	argInterpreter
	data   data.Data
	params params.Params
}

func (itp _interpreter) Interpret(
	inputArgs map[string]interface{},
	inputParams map[string]*model.Param,
	parentPkgHandle model.PkgHandle,
	pkgPath string,
	scope map[string]*model.Value,
	opScratchDir string,
) (map[string]*model.Value, []error) {
	interpretedArgs := map[string]*model.Value{}

	// 1) interpret
	var errs []error
	for argName, argValue := range inputArgs {
		var err error
		interpretedArgs[argName], err = itp.argInterpreter.Interpret(
			argName,
			argValue,
			inputParams[argName],
			parentPkgHandle,
			scope,
			opScratchDir,
		)
		if nil != err {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return nil, errs
	}

	// 2) apply defaults
	argsWithDefaults := itp.params.Default(interpretedArgs, inputParams, pkgPath)

	// 3) validate
	for inputName, inputValue := range argsWithDefaults {
		errs = append(errs, itp.data.Validate(inputValue, inputParams[inputName])...)
	}

	return argsWithDefaults, errs
}

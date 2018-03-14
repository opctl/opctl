package outputs

import (
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/opcall/params"
)

type interpreter interface {
	// Interpret applies defaults to & validates output args
	Interpret(
		outputArgs map[string]*model.Value,
		outputParams map[string]*model.Param,
		pkgPath string,
	) (
		map[string]*model.Value,
		error,
	)
}

func newInterpreter() interpreter {
	return _interpreter{
		params:    params.New(),
		validator: newValidator(),
	}
}

type _interpreter struct {
	params params.Params
	validator
}

func (itp _interpreter) Interpret(
	outputArgs map[string]*model.Value,
	outputParams map[string]*model.Param,
	pkgPath string,
) (
	map[string]*model.Value,
	error,
) {

	argsWithDefaults := itp.params.Default(outputArgs, outputParams, pkgPath)

	return argsWithDefaults, itp.validator.Validate(argsWithDefaults, outputParams)

}

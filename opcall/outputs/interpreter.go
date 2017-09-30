package outputs

import (
	"github.com/opspec-io/sdk-golang/model"
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
		defaulter: newDefaulter(),
		validator: newValidator(),
	}
}

type _interpreter struct {
	defaulter
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

	outputArgsWithDefaults := itp.defaulter.Default(outputArgs, outputParams, pkgPath)

	return outputArgsWithDefaults, itp.validator.Validate(outputArgsWithDefaults, outputParams)

}

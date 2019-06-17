package outputs

//go:generate counterfeiter -o ./fakeInterpreter.go --fake-name FakeInterpreter ./ Interpreter

import (
	"github.com/opctl/sdk-golang/model"
	"github.com/opctl/sdk-golang/opspec/interpreter/call/op/params"
)

type Interpreter interface {
	// Interpret applies defaults to & validates output args
	Interpret(
		outputArgs map[string]*model.Value,
		outputParams map[string]*model.Param,
		opPath string,
	) (
		map[string]*model.Value,
		error,
	)
}

// NewInterpreter returns an initialized Interpreter instance
func NewInterpreter() Interpreter {
	return _interpreter{
		paramsDefaulter: params.NewDefaulter(),
		paramsValidator: params.NewValidator(),
	}
}

type _interpreter struct {
	paramsDefaulter params.Defaulter
	paramsValidator params.Validator
}

func (itp _interpreter) Interpret(
	outputArgs map[string]*model.Value,
	outputParams map[string]*model.Param,
	opPath string,
) (
	map[string]*model.Value,
	error,
) {

	argsWithDefaults := itp.paramsDefaulter.Default(outputArgs, outputParams, opPath)

	return argsWithDefaults, itp.paramsValidator.Validate(argsWithDefaults, outputParams)

}

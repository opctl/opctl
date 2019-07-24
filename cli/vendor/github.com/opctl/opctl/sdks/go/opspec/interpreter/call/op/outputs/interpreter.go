package outputs

//go:generate counterfeiter -o ./fakeInterpreter.go --fake-name FakeInterpreter ./ Interpreter

import (
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/op/params"
	"github.com/opctl/opctl/sdks/go/types"
)

type Interpreter interface {
	// Interpret applies defaults to & validates output args
	Interpret(
		outputArgs map[string]*types.Value,
		outputParams map[string]*types.Param,
		opPath string,
	) (
		map[string]*types.Value,
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
	outputArgs map[string]*types.Value,
	outputParams map[string]*types.Param,
	opPath string,
) (
	map[string]*types.Value,
	error,
) {

	argsWithDefaults := itp.paramsDefaulter.Default(outputArgs, outputParams, opPath)

	return argsWithDefaults, itp.paramsValidator.Validate(argsWithDefaults, outputParams)

}

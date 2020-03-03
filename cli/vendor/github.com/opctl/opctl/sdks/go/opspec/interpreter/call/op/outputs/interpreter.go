package outputs

import (
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/op/params"
)

//counterfeiter:generate -o fakes/interpreter.go . Interpreter
type Interpreter interface {
	// Interpret applies defaults to & validates output args
	Interpret(
		outputArgs map[string]*model.Value,
		outputParams map[string]*model.Param,
		opPath string,
		opScratchDir string,
	) (
		map[string]*model.Value,
		error,
	)
}

// NewInterpreter returns an initialized Interpreter instance
func NewInterpreter() Interpreter {
	return _interpreter{
		paramsCoercer:   params.NewCoercer(),
		paramsDefaulter: params.NewDefaulter(),
		paramsValidator: params.NewValidator(),
	}
}

type _interpreter struct {
	paramsCoercer   params.Coercer
	paramsDefaulter params.Defaulter
	paramsValidator params.Validator
}

func (itp _interpreter) Interpret(
	outputArgs map[string]*model.Value,
	outputParams map[string]*model.Param,
	opPath string,
	opScratchDir string,
) (
	map[string]*model.Value,
	error,
) {

	var err error
	outputArgs, err = itp.paramsCoercer.Coerce(outputArgs, outputParams, opScratchDir)
	if nil != err {
		return outputArgs, err
	}

	argsWithDefaults := itp.paramsDefaulter.Default(outputArgs, outputParams, opPath)

	return argsWithDefaults, itp.paramsValidator.Validate(argsWithDefaults, outputParams)

}

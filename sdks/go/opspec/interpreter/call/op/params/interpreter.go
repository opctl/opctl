package params

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -o ./fakeInterpreter.go --fake-name FakeInterpreter ./ Interpreter

import (
	"github.com/opctl/opctl/sdks/go/model"
)

type Interpreter interface {
	// Interpret applies defaults to & validates args
	Interpret(
		scope map[string]*model.Value,
		params map[string]*model.Param,
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
		paramsCoercer:   NewCoercer(),
		paramsDefaulter: NewDefaulter(),
		paramsValidator: NewValidator(),
	}
}

type _interpreter struct {
	paramsCoercer   Coercer
	paramsDefaulter Defaulter
	paramsValidator Validator
}

func (itp _interpreter) Interpret(
	scope map[string]*model.Value,
	params map[string]*model.Param,
	opPath string,
	opScratchDir string,
) (
	map[string]*model.Value,
	error,
) {

	var err error
	scope, err = itp.paramsCoercer.Coerce(scope, params, opScratchDir)
	if nil != err {
		return scope, err
	}

	argsWithDefaults := itp.paramsDefaulter.Default(scope, params, opPath)

	return argsWithDefaults, itp.paramsValidator.Validate(argsWithDefaults, params)

}

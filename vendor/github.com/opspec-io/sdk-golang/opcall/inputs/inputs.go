package inputs

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ Inputs

import (
	"github.com/opspec-io/sdk-golang/model"
)

type Inputs interface {
	// Validate validates inputVals against inputParams
	Validate(
		inputVals map[string]*model.Data,
		inputParams map[string]*model.Param,
	) map[string][]error

	// Interpret interprets inputs via the provided inputArgs, inputParams, and scope
	Interpret(
		inputArgs map[string]string,
		inputParams map[string]*model.Param,
		scope map[string]*model.Data,
	) (map[string]*model.Data, []error)
}

func New() Inputs {
	return _Inputs{
		argInterpreter: newArgInterpreter(),
		validator:      newValidator(),
	}
}

type _Inputs struct {
	argInterpreter argInterpreter
	validator      validator
}

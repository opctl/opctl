package inputs

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ Inputs

import (
	"github.com/opspec-io/sdk-golang/data"
	"github.com/opspec-io/sdk-golang/model"
)

type Inputs interface {
	// Validate validates inputVals against inputParams
	// note: param defaults aren't considered
	Validate(
		inputVals map[string]*model.Value,
		inputParams map[string]*model.Param,
	) map[string][]error

	interpreter
}

func New() Inputs {
	return _Inputs{
		argInterpreter: newArgInterpreter(),
		data:           data.New(),
		interpreter:    newInterpreter(),
	}
}

type _Inputs struct {
	argInterpreter argInterpreter
	data           data.Data
	interpreter
}

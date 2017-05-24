package interpreter

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ Interpreter

import (
	"github.com/opspec-io/sdk-golang/interpolater"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/opcall/input/validator"
)

type Interpreter interface {
	// Interpret interprets an SCGOpCall input into an DCGOpCall input
	Interpret(
		name,
		value string,
		params map[string]*model.Param,
		scope map[string]*model.Data,
	) (*model.Data, error)
}

func New() Interpreter {
	return _Interpreter{
		interpolater: interpolater.New(),
		validator:    validator.New(),
	}
}

type _Interpreter struct {
	interpolater interpolater.Interpolater
	validator    validator.Validator
}

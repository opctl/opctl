package string

//go:generate counterfeiter -o ./fakeInterpreter.go --fake-name fakeInterpreter ./ Interpreter

import (
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/util/interpolater"
)

type Interpreter interface {
	// interprets an expression to a string
	Interpret(
		scope map[string]*model.Value,
		expression string,
	) (string, error)
}

func newInterpreter() Interpreter {
	return _interpreter{
		deReferencerFactory: newDeReferencerFactory(),
		interpolater:        interpolater.New(),
	}
}

type _interpreter struct {
	deReferencerFactory deReferencerFactory
	interpolater        interpolater.Interpolater
}

func (itp _interpreter) Interpret(
	scope map[string]*model.Value,
	expression string,
) (string, error) {
	return itp.interpolater.Interpolate(
		expression,
		itp.deReferencerFactory.New(scope),
	)
}

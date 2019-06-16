package loopable

import (
	"github.com/opctl/opctl/sdk/go/model"
	"github.com/opctl/opctl/sdk/go/opspec/interpreter/array"
	"github.com/opctl/opctl/sdk/go/opspec/interpreter/object"
)

//go:generate counterfeiter -o ./fakeInterpreter.go --fake-name FakeInterpreter ./ Interpreter

type Interpreter interface {
	Interpret(
		expression interface{},
		opHandle model.DataHandle,
		scope map[string]*model.Value,
	) (*model.Value, error)
}

// NewInterpreter returns an initialized Interpreter instance
func NewInterpreter() Interpreter {
	return &_interpreter{
		arrayInterpreter:  array.NewInterpreter(),
		objectInterpreter: object.NewInterpreter(),
	}
}

type _interpreter struct {
	arrayInterpreter  array.Interpreter
	objectInterpreter object.Interpreter
}

func (itp _interpreter) Interpret(
	expression interface{},
	opHandle model.DataHandle,
	scope map[string]*model.Value,
) (*model.Value, error) {
	// try interpreting as array
	if dcgForEach, err := itp.arrayInterpreter.Interpret(
		scope,
		expression,
		opHandle,
	); nil == err {
		return dcgForEach, err
	}

	// fallback to interpreting as object
	return itp.objectInterpreter.Interpret(
		scope,
		expression,
		opHandle,
	)
}

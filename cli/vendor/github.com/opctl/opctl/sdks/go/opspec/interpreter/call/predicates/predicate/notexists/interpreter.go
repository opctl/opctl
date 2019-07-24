package notexists

import (
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/reference"
	"github.com/opctl/opctl/sdks/go/types"
)

//go:generate counterfeiter -o ./fakeInterpreter.go --fake-name FakeInterpreter ./ Interpreter

type Interpreter interface {
	Interpret(
		expression string,
		opHandle types.DataHandle,
		scope map[string]*types.Value,
	) (bool, error)
}

// NewInterpreter returns an initialized Interpreter instance
func NewInterpreter() Interpreter {
	return &_interpreter{
		reference: reference.NewInterpreter(),
	}
}

type _interpreter struct {
	reference reference.Interpreter
}

func (itp _interpreter) Interpret(
	expression string,
	opHandle types.DataHandle,
	scope map[string]*types.Value,
) (bool, error) {
	// @TODO: make more exact. reference.Interpret can err for more reasons than simply null pointer exceptions.
	_, err := itp.reference.Interpret(
		expression,
		scope,
		opHandle,
	)

	return nil != err, nil
}

package exists

import (
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/reference"
)

//counterfeiter:generate -o fakes/interpreter.go . Interpreter
type Interpreter interface {
	Interpret(
		expression string,
		scope map[string]*model.Value,
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
	scope map[string]*model.Value,
) (bool, error) {

	// @TODO: make more exact. reference.Interpret can err for more reasons than simply null pointer exceptions.
	_, err := itp.reference.Interpret(
		expression,
		scope,
	)

	return nil == err, nil
}

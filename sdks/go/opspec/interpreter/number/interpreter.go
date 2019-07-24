package number

//go:generate counterfeiter -o ./fakeInterpreter.go --fake-name FakeInterpreter ./ Interpreter

import (
	"fmt"
	"github.com/opctl/opctl/sdks/go/data/coerce"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/value"
	"github.com/opctl/opctl/sdks/go/types"
)

type Interpreter interface {
	// Interpret interprets an expression to a number value.
	// Expression must be either a type supported by coerce.ToNumber
	// or an number initializer
	Interpret(
		scope map[string]*types.Value,
		expression interface{},
		opHandle types.DataHandle,
	) (*types.Value, error)
}

// NewInterpreter returns an initialized Interpreter instance
func NewInterpreter() Interpreter {
	return _interpreter{
		coerce:           coerce.New(),
		valueInterpreter: value.NewInterpreter(),
	}
}

type _interpreter struct {
	coerce           coerce.Coerce
	valueInterpreter value.Interpreter
}

func (itp _interpreter) Interpret(
	scope map[string]*types.Value,
	expression interface{},
	opHandle types.DataHandle,
) (*types.Value, error) {
	value, err := itp.valueInterpreter.Interpret(
		expression,
		scope,
		opHandle,
	)
	if nil != err {
		return nil, fmt.Errorf("unable to interpret %+v to number; error was %v", expression, err)
	}

	return itp.coerce.ToNumber(&value)
}

package eq

import (
	"github.com/opctl/sdk-golang/model"
	stringPkg "github.com/opctl/sdk-golang/opspec/interpreter/string"
)

//go:generate counterfeiter -o ./fakeInterpreter.go --fake-name FakeInterpreter ./ Interpreter

type Interpreter interface {
	Interpret(
		expressions []interface{},
		opHandle model.DataHandle,
		scope map[string]*model.Value,
	) (bool, error)
}

// NewInterpreter returns an initialized Interpreter instance
func NewInterpreter() Interpreter {
	return &_interpreter{
		stringInterpreter: stringPkg.NewInterpreter(),
	}
}

type _interpreter struct {
	stringInterpreter stringPkg.Interpreter
}

func (itp _interpreter) Interpret(
	expressions []interface{},
	opHandle model.DataHandle,
	scope map[string]*model.Value,
) (bool, error) {
	var firstItemAsString string
	for i, expression := range expressions {
		// interpret items as strings since everything is coercible to string
		item, err := itp.stringInterpreter.Interpret(scope, expression, opHandle)
		if nil != err {
			return false, err
		}
		currentItemAsString := *item.String

		if i == 0 {
			firstItemAsString = currentItemAsString
			continue
		}

		if firstItemAsString != currentItemAsString {
			// if first seen item not equal to current item predicate is false.
			return false, nil
		}
	}
	return true, nil
}

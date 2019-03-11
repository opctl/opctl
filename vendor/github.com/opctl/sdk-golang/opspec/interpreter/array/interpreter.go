package array

//go:generate counterfeiter -o ./fakeInterpreter.go --fake-name FakeInterpreter ./ Interpreter

import (
	"fmt"
	"github.com/opctl/sdk-golang/data/coerce"
	"github.com/opctl/sdk-golang/model"
	"github.com/opctl/sdk-golang/opspec/interpreter"
	"github.com/opctl/sdk-golang/opspec/interpreter/array/initializer"
	"github.com/opctl/sdk-golang/opspec/interpreter/interpolater"
)

type Interpreter interface {
	// Interpret evaluates an expression to an array value.
	// Expression must be either a type supported by coerce.ToArray
	// or an array initializer
	Interpret(
		scope map[string]*model.Value,
		expression interface{},
		opHandle model.DataHandle,
	) (*model.Value, error)
}

// NewInterpreter returns an initialized Interpreter instance
func NewInterpreter() Interpreter {
	return _interpreter{
		initializerInterpreter: initializer.NewInterpreter(),
		coerce:                 coerce.New(),
		interpolater:           interpolater.New(),
	}
}

type _interpreter struct {
	initializerInterpreter initializer.Interpreter
	coerce                 coerce.Coerce
	interpolater           interpolater.Interpolater
}

func (ea _interpreter) Interpret(
	scope map[string]*model.Value,
	expression interface{},
	opHandle model.DataHandle,
) (*model.Value, error) {
	switch expression := expression.(type) {
	case []interface{}:
		arrayValue, err := ea.initializerInterpreter.Interpret(
			expression,
			scope,
			opHandle,
		)
		if nil != err {
			return nil, fmt.Errorf("unable to evaluate %+v to array; error was %v", expression, err)
		}

		return &model.Value{Array: arrayValue}, nil
	case string:
		var value *model.Value
		if ref, ok := interpreter.TryResolveExplicitRef(expression, scope); ok {
			value = ref
		} else {
			stringValue, err := ea.interpolater.Interpolate(
				expression,
				scope,
				opHandle,
			)
			if nil != err {
				return nil, err
			}
			value = &model.Value{String: &stringValue}
		}
		return ea.coerce.ToArray(value)
	}
	return nil, fmt.Errorf("unable to evaluate %+v to array; unsupported type", expression)
}

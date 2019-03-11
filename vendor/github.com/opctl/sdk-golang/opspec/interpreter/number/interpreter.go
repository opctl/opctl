package number

//go:generate counterfeiter -o ./fakeInterpreter.go --fake-name FakeInterpreter ./ Interpreter

import (
	"fmt"
	"github.com/opctl/sdk-golang/data/coerce"
	"github.com/opctl/sdk-golang/model"
	"github.com/opctl/sdk-golang/opspec/interpreter"
	"github.com/opctl/sdk-golang/opspec/interpreter/interpolater"
)

type Interpreter interface {
	// Interpret evaluates an expression to a number value.
	// Expression must be a type supported by coerce.ToNumber
	Interpret(
		scope map[string]*model.Value,
		expression interface{},
		opHandle model.DataHandle,
	) (*model.Value, error)
}

// NewInterpreter returns an initialized Interpreter instance
func NewInterpreter() Interpreter {
	return _interpreter{
		coerce:       coerce.New(),
		interpolater: interpolater.New(),
	}
}

type _interpreter struct {
	coerce       coerce.Coerce
	interpolater interpolater.Interpolater
}

func (en _interpreter) Interpret(
	scope map[string]*model.Value,
	expression interface{},
	opHandle model.DataHandle,
) (*model.Value, error) {
	switch expression := expression.(type) {
	case float64:
		return &model.Value{Number: &expression}, nil
	case int:
		expAsFloat64 := float64(expression)
		return &model.Value{Number: &expAsFloat64}, nil
	case string:
		var value *model.Value
		if ref, ok := interpreter.TryResolveExplicitRef(expression, scope); ok {
			value = ref
		} else {
			stringValue, err := en.interpolater.Interpolate(
				expression,
				scope,
				opHandle,
			)
			if nil != err {
				return nil, err
			}

			value = &model.Value{String: &stringValue}
		}
		return en.coerce.ToNumber(value)
	}

	return nil, fmt.Errorf("unable to evaluate %+v to number; unsupported type", expression)
}

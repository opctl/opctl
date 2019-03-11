package object

//go:generate counterfeiter -o ./fakeInterpreter.go --fake-name FakeInterpreter ./ Interpreter

import (
	"fmt"
	"github.com/opctl/sdk-golang/data/coerce"
	"github.com/opctl/sdk-golang/model"
	"github.com/opctl/sdk-golang/opspec/interpreter"
	"github.com/opctl/sdk-golang/opspec/interpreter/interpolater"
	"github.com/opctl/sdk-golang/opspec/interpreter/object/initializer"
)

type Interpreter interface {
	// Interpret interprets an expression to a object value.
	// Expression must be either a type supported by coerce.ToObject
	// or an object initializer
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

func (eo _interpreter) Interpret(
	scope map[string]*model.Value,
	expression interface{},
	opHandle model.DataHandle,
) (*model.Value, error) {
	switch expression := expression.(type) {
	case map[string]interface{}:
		objectValue, err := eo.initializerInterpreter.Interpret(
			expression,
			scope,
			opHandle,
		)
		if nil != err {
			return nil, fmt.Errorf("unable to interpretuate %+v to object; error was %v", expression, err)
		}

		return &model.Value{Object: objectValue}, nil
	case string:
		var value *model.Value
		if ref, ok := interpreter.TryResolveExplicitRef(expression, scope); ok {
			value = ref
		} else {
			stringValue, err := eo.interpolater.Interpolate(
				expression,
				scope,
				opHandle,
			)
			if nil != err {
				return nil, err
			}
			value = &model.Value{String: &stringValue}
		}
		return eo.coerce.ToObject(value)
	}

	return nil, fmt.Errorf("unable to interpretuate %+v to object; unsupported type", expression)
}

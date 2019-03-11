package string

//go:generate counterfeiter -o ./fakeInterpreter.go --fake-name FakeInterpreter ./ Interpreter

import (
	"fmt"

	"github.com/opctl/sdk-golang/data/coerce"
	"github.com/opctl/sdk-golang/model"
	"github.com/opctl/sdk-golang/opspec/interpreter"
	arrayInitializer "github.com/opctl/sdk-golang/opspec/interpreter/array/initializer"
	"github.com/opctl/sdk-golang/opspec/interpreter/interpolater"
	objectInitializer "github.com/opctl/sdk-golang/opspec/interpreter/object/initializer"
)

type Interpreter interface {
	// Interpret evaluates an expression to a string value.
	Interpret(
		scope map[string]*model.Value,
		expression interface{},
		opHandle model.DataHandle,
	) (*model.Value, error)
}

// NewInterpreter returns an initialized Interpreter instance
func NewInterpreter() Interpreter {
	return _interpreter{
		arrayInitializerInterpreter:  arrayInitializer.NewInterpreter(),
		objectInitializerInterpreter: objectInitializer.NewInterpreter(),
		coerce:       coerce.New(),
		interpolater: interpolater.New(),
	}
}

type _interpreter struct {
	arrayInitializerInterpreter  arrayInitializer.Interpreter
	objectInitializerInterpreter objectInitializer.Interpreter
	coerce                       coerce.Coerce
	interpolater                 interpolater.Interpolater
}

func (es _interpreter) Interpret(
	scope map[string]*model.Value,
	expression interface{},
	opHandle model.DataHandle,
) (*model.Value, error) {
	var value *model.Value

	switch expression := expression.(type) {
	case bool:
		value = &model.Value{Boolean: &expression}
	case float64:
		value = &model.Value{Number: &expression}
	case int:
		expressionAsFloat64 := float64(expression)
		value = &model.Value{Number: &expressionAsFloat64}
	case map[string]interface{}:
		objectValue, err := es.objectInitializerInterpreter.Interpret(
			expression,
			scope,
			opHandle,
		)
		if nil != err {
			return nil, fmt.Errorf("unable to evaluate %+v to string; error was %v", expression, err)
		}

		value = &model.Value{Object: objectValue}
	case []interface{}:
		arrayValue, err := es.arrayInitializerInterpreter.Interpret(
			expression,
			scope,
			opHandle,
		)
		if nil != err {
			return nil, fmt.Errorf("unable to evaluate %+v to string; error was %v", expression, err)
		}

		value = &model.Value{Array: arrayValue}
	case string:
		if ref, ok := interpreter.TryResolveExplicitRef(expression, scope); ok {
			value = ref
		} else {
			stringValue, err := es.interpolater.Interpolate(
				expression,
				scope,
				opHandle,
			)
			if nil != err {
				return nil, err
			}

			value = &model.Value{String: &stringValue}
		}
	default:
		return nil, fmt.Errorf("unable to evaluate %+v to string; unsupported type", expression)
	}

	return es.coerce.ToString(value)
}

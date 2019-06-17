package value

//go:generate counterfeiter -o ./fakeInterpreter.go --fake-name FakeInterpreter ./ Interpreter

import (
	"fmt"

	"github.com/opctl/sdk-golang/model"
	"github.com/opctl/sdk-golang/opspec/interpreter/interpolater"
)

type Interpreter interface {
	// Interpret interprets a value expression
	Interpret(
		valueExpression interface{},
		scope map[string]*model.Value,
		opHandle model.DataHandle,
	) (model.Value, error)
}

// NewInterpreter returns an initialized Interpreter instance
func NewInterpreter() Interpreter {
	return _interpreter{
		interpolater: interpolater.New(),
	}
}

type _interpreter struct {
	interpolater interpolater.Interpolater
}

func (itp _interpreter) Interpret(
	valueExpression interface{},
	scope map[string]*model.Value,
	opHandle model.DataHandle,
) (model.Value, error) {
	switch typedValueExpression := valueExpression.(type) {
	case bool:
		return model.Value{Boolean: &typedValueExpression}, nil
	case float64:
		return model.Value{Number: &typedValueExpression}, nil
	case int:
		number := float64(typedValueExpression)
		return model.Value{Number: &number}, nil
	case map[string]interface{}:
		// object initializer
		value := map[string]interface{}{}
		for propertyKeyExpression, propertyValueExpression := range typedValueExpression {
			propertyKey, err := itp.interpolater.Interpolate(
				fmt.Sprintf("%v", propertyKeyExpression),
				scope,
				opHandle,
			)
			if nil != err {
				return model.Value{}, err
			}

			if nil == propertyValueExpression {
				// implicit reference
				propertyValueExpression = fmt.Sprintf("$(%v)", propertyKeyExpression)
			}
			boxedPropertyValue, err := itp.Interpret(
				propertyValueExpression,
				scope,
				opHandle,
			)
			if nil != err {
				return model.Value{}, fmt.Errorf("unable to interpret '%v: %v' as object initializer property; error was %v", propertyKeyExpression, propertyValueExpression, err)
			}

			if value[propertyKey], err = boxedPropertyValue.Unbox(); nil != err {
				return model.Value{}, fmt.Errorf("unable to interpret '%v: %v' as object initializer property; error was %v", propertyKeyExpression, propertyValueExpression, err)
			}
		}

		return model.Value{Object: &value}, nil
	case []interface{}:
		// array initializer
		value := []interface{}{}
		for _, itemExpression := range typedValueExpression {
			boxedItemValue, err := itp.Interpret(
				itemExpression,
				scope,
				opHandle,
			)
			if nil != err {
				return model.Value{}, fmt.Errorf("unable to interpret '%+v' as array initializer item; error was %v", itemExpression, err)
			}

			itemValue, err := boxedItemValue.Unbox()
			if nil != err {
				return model.Value{}, fmt.Errorf("unable to interpret '%+v' as array initializer item; error was %v", itemExpression, err)
			}
			value = append(value, itemValue)
		}
		return model.Value{Array: &value}, nil
	case string:
		value, err := itp.interpolater.Interpolate(
			typedValueExpression,
			scope,
			opHandle,
		)
		if nil != err {
			return model.Value{}, err
		}

		return model.Value{String: &value}, nil
	default:
		return model.Value{}, fmt.Errorf("unable to interpret %+v to string; unsupported type", typedValueExpression)
	}
}

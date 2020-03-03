package value

import (
	"fmt"
	"regexp"

	"github.com/golang-interfaces/iioutil"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/interpolater"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/reference"
)

//counterfeiter:generate -o fakes/interpreter.go . Interpreter
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
		interpolater:         interpolater.New(),
		ioUtil:               iioutil.New(),
		referenceInterpreter: reference.NewInterpreter(),
	}
}

type _interpreter struct {
	interpolater         interpolater.Interpolater
	ioUtil               iioutil.IIOUtil
	referenceInterpreter reference.Interpreter
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
				propertyKeyExpression,
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
			propertyValue, err := itp.Interpret(
				propertyValueExpression,
				scope,
				opHandle,
			)
			if nil != err {
				return model.Value{}, fmt.Errorf("unable to interpret '%v: %v' as object initializer property; error was %v", propertyKeyExpression, propertyValueExpression, err)
			}

			if nil != propertyValue.File {
				fileBytes, err := itp.ioUtil.ReadFile(*propertyValue.File)
				if nil != err {
					return model.Value{}, fmt.Errorf("unable to interpret '%v: %v' as object initializer property; error was %v", propertyKeyExpression, propertyValueExpression, err)
				}

				value[propertyKey] = string(fileBytes)
				continue
			} else if nil != propertyValue.Dir {
				return model.Value{}, fmt.Errorf("unable to interpret '%v: %v' as object initializer property; directories aren't valid object properties", propertyKeyExpression, propertyValueExpression)
			} else if nil != propertyValue.Socket {
				return model.Value{}, fmt.Errorf("unable to interpret '%v: %v' as object initializer property; sockets aren't valid object properties", propertyKeyExpression, propertyValueExpression)
			}

			unboxedPropertyValue, unboxErr := propertyValue.Unbox()
			if nil != unboxErr {
				return model.Value{}, unboxErr
			}

			value[propertyKey] = unboxedPropertyValue
		}

		return model.Value{Object: &value}, nil
	case []interface{}:
		// array initializer
		value := []interface{}{}
		for _, itemExpression := range typedValueExpression {
			itemValue, err := itp.Interpret(
				itemExpression,
				scope,
				opHandle,
			)
			if nil != err {
				return model.Value{}, fmt.Errorf("unable to interpret '%+v' as array initializer item; error was %v", itemExpression, err)
			}
			value = append(value, itemValue)
		}
		return model.Value{Array: &value}, nil
	case string:
		if regexp.MustCompile("^\\$\\(.+\\)$").MatchString(typedValueExpression) {
			// attempt to process as a reference since its reference like.
			// @TODO: make more exact. reference.Interpret can err for reasons beyond not being a reference.
			value, err := itp.referenceInterpreter.Interpret(
				typedValueExpression,
				scope,
				opHandle,
			)
			if nil == err {
				return *value, nil
			}
		}

		// fallback to processing as a string
		var valueString string
		valueString, err := itp.interpolater.Interpolate(
			typedValueExpression,
			scope,
			opHandle,
		)
		if nil != err {
			return model.Value{}, err
		}

		return model.Value{String: &valueString}, nil
	case model.Value:
		return typedValueExpression, nil
	default:
		return model.Value{}, fmt.Errorf("unable to interpret %+v as value; unsupported type", typedValueExpression)
	}
}

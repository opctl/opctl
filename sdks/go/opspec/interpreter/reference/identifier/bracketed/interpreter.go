package bracketed

import (
	"fmt"
	"strings"

	"github.com/opctl/opctl/sdks/go/opspec/interpreter/reference/identifier/bracketed/item"

	"github.com/opctl/opctl/sdks/go/opspec/interpreter/reference/identifier/value"

	"github.com/opctl/opctl/sdks/go/model"
)

// Interpreter interprets a bracketed identifier from ref by consuming from '[' up to & including the first ']'
// it's an error if ref doesn't start with '[' or contain ']'
// returns ref remainder, dereferenced data, and error if one occurred
type Interpreter interface {
	Interpret(
		ref string,
		data *model.Value,
	) (string, *model.Value, error)
}

func NewInterpreter() Interpreter {
	return _interpreter{
		coerceToArrayOrObjecter: newCoerceToArrayOrObjecter(),
		itemInterpreter:         item.NewInterpreter(),
		valueConstructor:        value.NewConstructor(),
	}
}

type _interpreter struct {
	coerceToArrayOrObjecter coerceToArrayOrObjecter
	itemInterpreter         item.Interpreter
	valueConstructor        value.Constructor
}

func (dr _interpreter) Interpret(
	ref string,
	data *model.Value,
) (string, *model.Value, error) {

	if !strings.HasPrefix(ref, "[") {
		return "", nil, fmt.Errorf("unable to interpret '%v'; expected '['", ref)
	}

	indexOfNextCloseBracket := strings.Index(ref, "]")
	if indexOfNextCloseBracket < 0 {
		return "", nil, fmt.Errorf("unable to interpret '%v'; expected ']'", ref)
	}

	data, err := dr.coerceToArrayOrObjecter.CoerceToArrayOrObject(data)
	if nil != err {
		return "", nil, fmt.Errorf("unable to interpret '%v'; error was %v", ref, err.Error())
	}

	identifier := ref[1:indexOfNextCloseBracket]
	refRemainder := ref[indexOfNextCloseBracket+1:]

	if nil != data.Array {
		// data is array
		itemValue, err := dr.itemInterpreter.Interpret(identifier, *data)
		if nil != err {
			return "", nil, err
		}

		return refRemainder, itemValue, nil
	}

	// data is object
	property := (*data.Object)[identifier]
	propertyValue, err := dr.valueConstructor.Construct(property)
	if nil != err {
		return "", nil, fmt.Errorf("unable to interpret property; error was %v", err.Error())
	}
	return refRemainder, propertyValue, nil
}

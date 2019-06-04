package unbracketed

import (
	"fmt"

	"github.com/opctl/sdk-golang/data/coerce"
	"github.com/opctl/sdk-golang/model"
	"github.com/opctl/sdk-golang/opspec/interpreter/reference/identifier/value"
)

// Interpreter interprets an unbracketed identifier from ref as determined by unbracketed/parser.go
// returns remaining ref, dereferenced data, and error if one occurred
type Interpreter interface {
	Interpret(
		ref string,
		data *model.Value,
	) (string, *model.Value, error)
}

func NewInterpreter() Interpreter {
	return _interpreter{
		coerce:           coerce.New(),
		parser:           NewParser(),
		valueConstructor: value.NewConstructor(),
	}
}

type _interpreter struct {
	coerce           coerce.Coerce
	parser           Parser
	valueConstructor value.Constructor
}

func (dr _interpreter) Interpret(
	ref string,
	data *model.Value,
) (string, *model.Value, error) {

	dataAsObject, err := dr.coerce.ToObject(data)
	if nil != err {
		return ref, nil, fmt.Errorf("unable to interpret '%v'; error was %v", ref, err.Error())
	}

	identifier, refRemainder := dr.parser.Parse(ref)

	scopeValue, isValueInScope := (*dataAsObject.Object)[identifier]
	if !isValueInScope {
		return ref, nil, fmt.Errorf("unable to interpret '%v'; '%v' doesn't exist", ref, identifier)
	}

	identifierValue, err := dr.valueConstructor.Construct(scopeValue)
	if nil != err {
		return ref, nil, fmt.Errorf("unable to interpret '%v'; error was %v", ref, err.Error())
	}

	return refRemainder, identifierValue, nil
}

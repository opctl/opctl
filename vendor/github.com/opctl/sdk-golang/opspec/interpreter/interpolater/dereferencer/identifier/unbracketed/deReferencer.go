package unbracketed

import (
	"fmt"

	"github.com/opctl/sdk-golang/data/coerce"
	"github.com/opctl/sdk-golang/model"
	"github.com/opctl/sdk-golang/opspec/interpreter/interpolater/dereferencer/identifier/value"
)

// DeReferencer dereferences an unbracketed identifier from ref as determined by unbracketed/parser.go
// returns remaining ref, dereferenced data, and error if one occurred
type DeReferencer interface {
	DeReference(
		ref string,
		data *model.Value,
	) (string, *model.Value, error)
}

func NewDeReferencer() DeReferencer {
	return _deReferencer{
		coerce:           coerce.New(),
		parser:           NewParser(),
		valueConstructor: value.NewConstructor(),
	}
}

type _deReferencer struct {
	coerce           coerce.Coerce
	parser           Parser
	valueConstructor value.Constructor
}

func (dr _deReferencer) DeReference(
	ref string,
	data *model.Value,
) (string, *model.Value, error) {

	dataAsObject, err := dr.coerce.ToObject(data)
	if nil != err {
		return ref, nil, fmt.Errorf("unable to deReference '%v'; error was %v", ref, err.Error())
	}

	identifier, refRemainder := dr.parser.Parse(ref)

	scopeValue, isValueInScope := dataAsObject.Object[identifier]
	if !isValueInScope {
		return ref, nil, fmt.Errorf("unable to deReference '%v'; '%v' doesn't exist", ref, identifier)
	}

	identifierValue, err := dr.valueConstructor.Construct(scopeValue)
	if nil != err {
		return ref, nil, fmt.Errorf("unable to deReference '%v'; error was %v", ref, err.Error())
	}

	return refRemainder, identifierValue, nil
}

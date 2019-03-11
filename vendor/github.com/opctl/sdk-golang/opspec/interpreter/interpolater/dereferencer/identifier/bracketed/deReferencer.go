package bracketed

import (
	"fmt"
	"github.com/opctl/sdk-golang/opspec/interpreter/interpolater/dereferencer/identifier/bracketed/item"
	"strings"

	"github.com/opctl/sdk-golang/opspec/interpreter/interpolater/dereferencer/identifier/value"

	"github.com/opctl/sdk-golang/model"
)

// DeReferencer dereferences a bracketed identifier from ref by consuming from '[' up to & including the first ']'
// it's an error if ref doesn't start with '[' or contain ']'
// returns ref remainder, dereferenced data, and error if one occurred
type DeReferencer interface {
	DeReference(
		ref string,
		data *model.Value,
	) (string, *model.Value, error)
}

func NewDeReferencer() DeReferencer {
	return _deReferencer{
		coerceToArrayOrObjecter: newCoerceToArrayOrObjecter(),
		itemDeReferencer:        item.NewDeReferencer(),
		valueConstructor:        value.NewConstructor(),
	}
}

type _deReferencer struct {
	coerceToArrayOrObjecter coerceToArrayOrObjecter
	itemDeReferencer        item.DeReferencer
	valueConstructor        value.Constructor
}

func (dr _deReferencer) DeReference(
	ref string,
	data *model.Value,
) (string, *model.Value, error) {

	if !strings.HasPrefix(ref, "[") {
		return "", nil, fmt.Errorf("unable to deReference '%v'; expected '['", ref)
	}

	indexOfNextCloseBracket := strings.Index(ref, "]")
	if indexOfNextCloseBracket < 0 {
		return "", nil, fmt.Errorf("unable to deReference '%v'; expected ']'", ref)
	}

	data, err := dr.coerceToArrayOrObjecter.CoerceToArrayOrObject(data)
	if nil != err {
		return "", nil, fmt.Errorf("unable to deReference '%v'; error was %v", ref, err.Error())
	}

	identifier := ref[1:indexOfNextCloseBracket]
	refRemainder := ref[indexOfNextCloseBracket+1:]

	if nil != data.Array {
		// data is array
		itemValue, err := dr.itemDeReferencer.DeReference(identifier, *data)
		if nil != err {
			return "", nil, err
		}

		return refRemainder, itemValue, nil
	}

	// data is object
	property := data.Object[identifier]
	propertyValue, err := dr.valueConstructor.Construct(property)
	if nil != err {
		return "", nil, fmt.Errorf("unable to deReference property; error was %v", err.Error())
	}
	return refRemainder, propertyValue, nil
}

package unbracketed

import (
	"fmt"

	"github.com/opctl/opctl/sdks/go/data/coerce"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/reference/identifier/value"
	"github.com/pkg/errors"
)

// Interpret an unbracketed identifier from ref as determined by unbracketed/parse.go
// returns remaining ref, dereferenced data, and error if one occurred
func Interpret(
	ref string,
	data *model.Value,
) (string, *model.Value, error) {

	dataAsObject, err := coerce.ToObject(data)
	if nil != err {
		return ref, nil, errors.Wrap(err, fmt.Sprintf("unable to interpret '%v'", ref))
	}

	identifier, refRemainder := Parse(ref)

	scopeValue, isValueInScope := (*dataAsObject.Object)[identifier]
	if !isValueInScope {
		return ref, nil, fmt.Errorf("unable to interpret '%v': '%v' doesn't exist", ref, identifier)
	}

	identifierValue, err := value.Construct(scopeValue)
	if nil != err {
		return ref, nil, errors.Wrap(err, fmt.Sprintf("unable to interpret '%v'", ref))
	}

	return refRemainder, identifierValue, nil
}

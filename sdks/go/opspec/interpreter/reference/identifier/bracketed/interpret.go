package bracketed

import (
	"fmt"
	"strings"

	"github.com/opctl/opctl/sdks/go/opspec/interpreter/reference/identifier/bracketed/item"
	"github.com/pkg/errors"

	"github.com/opctl/opctl/sdks/go/opspec/interpreter/reference/identifier/value"

	"github.com/opctl/opctl/sdks/go/model"
)

// Interpret a bracketed identifier from ref by consuming from '[' up to & including the first ']'
// it's an error if ref doesn't start with '[' or contain ']'
// returns ref remainder, dereferenced data, and error if one occurred
func Interpret(
	ref string,
	data *model.Value,
) (string, *model.Value, error) {

	if !strings.HasPrefix(ref, "[") {
		return "", nil, fmt.Errorf("unable to interpret '%v': expected '['", ref)
	}

	indexOfNextCloseBracket := strings.Index(ref, "]")
	if indexOfNextCloseBracket < 0 {
		return "", nil, fmt.Errorf("unable to interpret '%v': expected ']'", ref)
	}

	data, err := CoerceToArrayOrObject(data)
	if nil != err {
		return "", nil, errors.Wrap(err, fmt.Sprintf("unable to interpret '%v'", ref))
	}

	identifier := ref[1:indexOfNextCloseBracket]
	refRemainder := ref[indexOfNextCloseBracket+1:]

	if nil != data.Array {
		// data is array
		itemValue, err := item.Interpret(identifier, *data)
		if nil != err {
			return "", nil, err
		}

		return refRemainder, itemValue, nil
	}

	// data is object
	property := (*data.Object)[identifier]
	propertyValue, err := value.Construct(property)
	if nil != err {
		return "", nil, errors.Wrap(err, "unable to interpret property")
	}
	return refRemainder, propertyValue, nil
}

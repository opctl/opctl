package bracketed

import (
	"fmt"
	"strings"

	"github.com/opctl/opctl/sdks/go/opspec/interpreter/reference/identifier/bracketed/item"

	"github.com/opctl/opctl/sdks/go/opspec/interpreter/reference/identifier/value"

	"github.com/opctl/opctl/sdks/go/model"
)

// Interpret a bracketed identifier from ref by consuming from '[' up to & including the first ']'
// it's an error if ref doesn't start with '[' or contain ']'
// returns ref remainder, dereferenced data, and error if one occurred
func Interpret(
	ref string,
	data *ipld.Node,
) (string, *ipld.Node, error) {

	if !strings.HasPrefix(ref, "[") {
		return "", nil, fmt.Errorf("unable to interpret '%v': expected '['", ref)
	}

	indexOfNextCloseBracket := strings.Index(ref, "]")
	if indexOfNextCloseBracket < 0 {
		return "", nil, fmt.Errorf("unable to interpret '%v': expected ']'", ref)
	}

	data, err := CoerceToArrayOrObject(data)
	if err != nil {
		return "", nil, fmt.Errorf("unable to interpret '%v': %w", ref, err)
	}

	identifier := ref[1:indexOfNextCloseBracket]
	refRemainder := ref[indexOfNextCloseBracket+1:]

	if data.Array != nil {
		// data is array
		itemValue, err := item.Interpret(identifier, *data)
		if err != nil {
			return "", nil, err
		}

		return refRemainder, itemValue, nil
	}

	// data is object
	property := (*data.Object)[identifier]
	propertyValue, err := value.Construct(property)
	if err != nil {
		return "", nil, fmt.Errorf("unable to interpret property: %w", err)
	}
	return refRemainder, propertyValue, nil
}

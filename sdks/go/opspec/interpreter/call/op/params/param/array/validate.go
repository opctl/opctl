package array

import (
	"github.com/opctl/opctl/sdks/go/data/coerce"
	"github.com/opctl/opctl/sdks/go/internal/jsonschema"
	"github.com/opctl/opctl/sdks/go/model"
)

// Validate validates a value against a string parameter
func Validate(
	value *model.Value,
	constraints map[string]interface{},
) []error {
	valueAsArray, err := coerce.ToArray(value)
	if err != nil {
		return []error{err}
	}

	valueAsArrayUnboxed, err := valueAsArray.Unbox()
	if err != nil {
		return []error{
			err,
		}
	}

	return jsonschema.Validate(
		valueAsArrayUnboxed,
		constraints,
	)
}

package object

import (
	"github.com/opctl/opctl/sdks/go/data/coerce"
	"github.com/opctl/opctl/sdks/go/internal/jsonschema"
	"github.com/opctl/opctl/sdks/go/model"
)

// Validate validates a value against an object parameter
func Validate(
	value *model.Value,
	constraints map[string]interface{},
) []error {

	valueAsObject, err := coerce.ToObject(value)
	if err != nil {
		return []error{err}
	}

	valueAsObjectUnboxed, err := valueAsObject.Unbox()
	if err != nil {
		return []error{
			err,
		}
	}
	
	return jsonschema.Validate(
		valueAsObjectUnboxed,
		constraints,
	)
}

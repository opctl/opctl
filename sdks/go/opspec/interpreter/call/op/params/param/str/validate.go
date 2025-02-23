package str

import (
	"github.com/opctl/opctl/sdks/go/data/coerce"
	"github.com/opctl/opctl/sdks/go/internal/jsonschema"
	"github.com/opctl/opctl/sdks/go/model"
)

// Validate validates a value against a number parameter
func Validate(
	value *model.Value,
	constraints map[string]interface{},
) []error {
	valueAsString, err := coerce.ToString(value)
	if err != nil {
		return []error{err}
	}

	return jsonschema.Validate(
		*valueAsString.String,
		constraints,
	)
}

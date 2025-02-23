package jsonschema

import (
	"github.com/santhosh-tekuri/jsonschema/v6"
)

// Validate validates a value against constraints
func Validate(
	value any,
	schema map[string]interface{},
) []error {
	// guard no schema
	if schema == nil {
		return nil
	}

	schemaID := ""
	compiler := jsonschema.NewCompiler()
	if err := compiler.AddResource(schemaID, schema); err != nil {
		return []error{err}
	}

	compiledSchema, err := compiler.Compile(schemaID)
	if err != nil {
		return []error{err}
	}

	if err := compiledSchema.Validate(value); err != nil {
		if vErr, ok := err.(*jsonschema.ValidationError); ok {
			var errs []error
			for _, cause := range vErr.Causes {
				errs = append(
					errs,
					cause,
				)
			}
			return errs
		}
		return []error{
			err,
		}
	}

	return nil
}

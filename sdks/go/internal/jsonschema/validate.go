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

	// register custom format checkers
	// gojsonschema.FormatCheckers.Add("docker-image-ref", formats.DockerImageRefFormatChecker{})
	// gojsonschema.FormatCheckers.Add("integer", formats.IntegerFormatChecker{})
	// gojsonschema.FormatCheckers.Add("semver", formats.SemVerFormatChecker{})
	schemaID := "schema.json" // Arbitrary ID for the schema
	compiler := jsonschema.NewCompiler()

	if err := compiler.AddResource(schemaID, schema); err != nil {
		return []error{err}
	}

	s, err := compiler.Compile(schemaID)
	if err != nil {
		return []error{err}
	}

	errorMessages := map[string]string{}
	extractErrorMessages(schema, "", errorMessages)

	// Validate the parsed document
	if err := s.Validate(value); err != nil {
		if vErr, ok := err.(*jsonschema.ValidationError); ok {
			var errs []error
			for _, cause := range vErr.Causes {
				errs = append(errs, cause)
			}
			return errs
		}
		return []error{err}
	}

	return nil
}

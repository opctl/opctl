package opfile

import (
	"bytes"
	"github.com/ghodss/yaml"
	"github.com/opctl/opctl/opspec/opfile"
	"github.com/santhosh-tekuri/jsonschema/v6"
)

// Validate validates an "op.yml"
func Validate(
	opFileBytes []byte,
) []error {
	schemaID := "op.yml"
	doc, err := jsonschema.UnmarshalJSON(bytes.NewReader(opfile.JsonSchemaBytes))
	if err != nil {
		return []error{err}
	}

	compiler := jsonschema.NewCompiler()
	if err := compiler.AddResource(schemaID, doc); err != nil {
		return []error{err}
	}

	opFileSchema, err := compiler.Compile(schemaID)
	if err != nil {
		return []error{err}
	}

	var unmarshalledYAML map[string]interface{}
	err = yaml.Unmarshal(opFileBytes, &unmarshalledYAML)
	if err != nil {
		// handle syntax errors specially
		return []error{err}
	}

	if err := opFileSchema.Validate(unmarshalledYAML); err != nil {
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

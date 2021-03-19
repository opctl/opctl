package opfile

//go:generate go run github.com/mjibson/esc -pkg=opfile -o validate_schema.go -private ../../../../opspec/opfile/jsonschema.json

import (
	"fmt"

	"github.com/ghodss/yaml"
	"github.com/xeipuuv/gojsonschema"
)

// Validate validates an "op.yml"
func Validate(
	opFileBytes []byte,
) []error {

	opFileSchemaBytes, err := _escFSByte(false, "/opspec/opfile/jsonschema.json")
	if err != nil {
		return []error{err}
	}

	opFileSchema, err := gojsonschema.NewSchema(
		gojsonschema.NewBytesLoader(opFileSchemaBytes),
	)
	if err != nil {
		return []error{err}
	}

	var unmarshalledYAML map[string]interface{}
	err = yaml.Unmarshal(opFileBytes, &unmarshalledYAML)
	if err != nil {
		// handle syntax errors specially
		return []error{err}
	}

	var errs []error
	result, err := opFileSchema.Validate(
		gojsonschema.NewGoLoader(unmarshalledYAML),
	)
	if err != nil {
		// handle syntax errors specially
		return append(errs, err)
	}
	for _, desc := range result.Errors() {
		errs = append(errs, fmt.Errorf("%s", desc))
	}

	return errs
}

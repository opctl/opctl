package opfile

//go:generate go run github.com/mjibson/esc -pkg=opfile -o validator_schema.go -private ../../../../opspec/opfile/jsonschema.json

import (
	"fmt"

	"github.com/ghodss/yaml"
	"github.com/xeipuuv/gojsonschema"
)

//counterfeiter:generate -o internal/fakes/validator.go . validator
type validator interface {
	// Validate validates an "op.yml"
	Validate(
		opFileBytes []byte,
	) []error
}

func newValidator() validator {
	opFileSchemaBytes, err := _escFSByte(false, "/opspec/opfile/jsonschema.json")
	if nil != err {
		panic(err)
	}

	opFileSchema, err := gojsonschema.NewSchema(
		gojsonschema.NewBytesLoader(opFileSchemaBytes),
	)
	if err != nil {
		panic(err)
	}

	return _validator{
		opFileSchema: opFileSchema,
	}
}

type _validator struct {
	opFileSchema *gojsonschema.Schema
}

func (vdr _validator) Validate(
	opFileBytes []byte,
) []error {

	var unmarshalledYAML map[string]interface{}
	err := yaml.Unmarshal(opFileBytes, &unmarshalledYAML)
	if nil != err {
		// handle syntax errors specially
		return []error{err}
	}

	var errs []error
	result, err := vdr.opFileSchema.Validate(
		gojsonschema.NewGoLoader(unmarshalledYAML),
	)
	if nil != err {
		// handle syntax errors specially
		return append(errs, err)
	}
	for _, desc := range result.Errors() {
		errs = append(errs, fmt.Errorf("%s", desc))
	}

	return errs
}

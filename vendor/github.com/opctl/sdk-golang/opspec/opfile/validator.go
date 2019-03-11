package dotyml

//go:generate counterfeiter -o ./fakeValidator.go --fake-name fakeValidator ./ validator

import (
	"fmt"
	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v2"
)

type validator interface {
	// Validate validates an "op.yml"
	Validate(
		manifestBytes []byte,
	) []error
}

func newValidator() validator {
	manifestSchemaBytes, err := githubComOpctlSpecsOpspecOpfileJsonschemaJsonBytes()
	if nil != err {
		panic(err)
	}

	manifestSchema, err := gojsonschema.NewSchema(
		gojsonschema.NewBytesLoader(manifestSchemaBytes),
	)
	if err != nil {
		panic(err)
	}

	return _validator{
		manifestSchema: manifestSchema,
	}
}

type _validator struct {
	manifestSchema *gojsonschema.Schema
}

func (vdr _validator) Validate(
	manifestBytes []byte,
) []error {

	unmarshalledYAML := map[string]interface{}{}
	err := yaml.Unmarshal(manifestBytes, unmarshalledYAML)
	if nil != err {
		// handle syntax errors specially
		return []error{err}
	}

	var errs []error
	result, err := vdr.manifestSchema.Validate(
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

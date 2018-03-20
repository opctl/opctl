package dotyml

//go:generate counterfeiter -o ./fakeValidator.go --fake-name fakeValidator ./ validator

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/xeipuuv/gojsonschema"
)

type validator interface {
	// Validate validates an "op.yml"
	Validate(
		manifestBytes []byte,
	) []error
}

func newValidator() validator {
	manifestSchemaBytes, err := githubComOpspecIoSpecSpecOpYmlSchemaJsonBytes()
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

	manifestJSONBytes, err := yaml.YAMLToJSON(manifestBytes)
	if nil != err {
		// handle syntax errors specially
		return []error{err}
	}

	var errs []error
	result, err := vdr.manifestSchema.Validate(
		gojsonschema.NewBytesLoader(manifestJSONBytes),
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

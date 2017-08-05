package manifest

//go:generate counterfeiter -o ./fakeValidator.go --fake-name fakeValidator ./ Validator

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/golang-interfaces/iioutil"
	"github.com/xeipuuv/gojsonschema"
)

type Validator interface {
	// Validate validates the pkg manifest
	Validate(
		manifestBytes []byte,
	) []error
}

func newValidator() Validator {

	// register custom format checkers
	gojsonschema.FormatCheckers.Add("uri-reference", uriRefFormatChecker{})

	manifestSchemaBytes, err := pkgManifestDataPkgManifestSchemaJsonBytes()
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
		ioUtil:         iioutil.New(),
		manifestSchema: manifestSchema,
	}
}

type _validator struct {
	ioUtil         iioutil.IIOUtil
	manifestSchema *gojsonschema.Schema
}

func (this _validator) Validate(
	manifestBytes []byte,
) []error {

	manifestJSONBytes, err := yaml.YAMLToJSON(manifestBytes)
	if nil != err {
		// handle syntax errors specially
		return []error{err}
	}

	var errs []error
	result, err := this.manifestSchema.Validate(
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

package manifest

//go:generate counterfeiter -o ./fakeValidator.go --fake-name fakeValidator ./ validator

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/golang-interfaces/iioutil"
	"github.com/xeipuuv/gojsonschema"
)

type validator interface {
	// Validate validates the pkg manifest at path
	Validate(
		path string,
	) []error
}

func newValidator() validator {

	// register custom format checkers
	gojsonschema.FormatCheckers.Add("uri-reference", uriRefFormatChecker{})

	manifestSchemaBytes, err := pkgManifestDataPackageManifestSchemaJsonBytes()
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
	ioUtil         iioutil.Iioutil
	manifestSchema *gojsonschema.Schema
}

func (this _validator) Validate(
	path string,
) []error {

	ManifestYAMLBytes, err := this.ioUtil.ReadFile(
		path,
	)
	if nil != err {
		// handle syntax errors specially
		return []error{err}
	}

	manifestJSONBytes, err := yaml.YAMLToJSON(ManifestYAMLBytes)
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

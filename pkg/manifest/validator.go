package manifest

//go:generate counterfeiter -o ./fakeValidator.go --fake-name fakeValidator ./ Validator

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/golang-interfaces/iioutil"
	"github.com/xeipuuv/gojsonschema"
	"io"
)

type Validator interface {
	// Validate validates the pkg manifest at path
	Validate(
		manifestReader io.Reader,
	) []error
}

func newValidator() Validator {

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
	ioUtil         iioutil.IIOUtil
	manifestSchema *gojsonschema.Schema
}

func (this _validator) Validate(
	manifestReader io.Reader,
) []error {

	ManifestYAMLBytes, err := this.ioUtil.ReadAll(manifestReader)
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

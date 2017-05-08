package pkg

//go:generate counterfeiter -o ./fakeManifestValidator.go --fake-name fakeManifestValidator ./ manifestValidator

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/golang-interfaces/iioutil"
	"github.com/xeipuuv/gojsonschema"
	"path"
)

type manifestValidator interface {
	Validate(pkgPath string) []error
}

func newManifestValidator() manifestValidator {
	manifestSchemaBytes, err := pkgDataPackageManifestSchemaJsonBytes()
	if nil != err {
		panic(err)
	}

	manifestSchema, err := gojsonschema.NewSchema(
		gojsonschema.NewBytesLoader(manifestSchemaBytes),
	)
	if err != nil {
		panic(err)
	}

	return _manifestValidator{
		ioUtil:         iioutil.New(),
		manifestSchema: manifestSchema,
	}

}

type _manifestValidator struct {
	ioUtil         iioutil.Iioutil
	manifestSchema *gojsonschema.Schema
}

func (this _manifestValidator) Validate(
	pkgPath string,
) []error {

	ManifestYAMLBytes, err := this.ioUtil.ReadFile(
		path.Join(pkgPath, OpDotYmlFileName),
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

package pkg

//go:generate counterfeiter -o ./fakeValidator.go --fake-name fakeValidator ./ validator

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/virtual-go/fs"
	"github.com/virtual-go/vioutil"
	"github.com/xeipuuv/gojsonschema"
	"path"
)

type validator interface {
	Validate(pkgRef string) []error
}

func newValidator(
	fs fs.FS,
) validator {
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

	return _validator{
		ioUtil:         vioutil.New(fs),
		manifestSchema: manifestSchema,
	}

}

type _validator struct {
	ioUtil         vioutil.VIOUtil
	manifestSchema *gojsonschema.Schema
}

func (this _validator) Validate(
	pkgRef string,
) []error {
	var errs []error

	ManifestYAMLBytes, err := this.ioUtil.ReadFile(
		path.Join(pkgRef, OpDotYmlFileName),
	)
	if nil != err {
		// handle syntax errors specially
		return append(errs, err)
	}

	manifestJSONBytes, err := yaml.YAMLToJSON(ManifestYAMLBytes)
	if nil != err {
		// handle syntax errors specially
		return append(errs, err)
	}

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

package pkg

//go:generate counterfeiter -o ./fakeValidator.go --fake-name fakeValidator ./ validator

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/opspec-io/sdk-golang/util/fs"
	"github.com/xeipuuv/gojsonschema"
	"path"
)

type validator interface {
	Validate(pkgRef string) (errs []error)
}

func newValidator() validator {
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
		fileSystem:     fs.NewFileSystem(),
		manifestSchema: manifestSchema,
	}

}

type _validator struct {
	fileSystem     fs.FileSystem
	manifestSchema *gojsonschema.Schema
}

func (this _validator) Validate(
	pkgRef string,
) (errs []error) {
	ManifestYAMLBytes, err := this.fileSystem.GetBytesOfFile(
		path.Join(pkgRef, NameOfPackageManifestFile),
	)
	if nil != err {
		// handle syntax errors specially
		errs = append(errs, fmt.Errorf("Error validating pkg. Details: %v", err.Error()))
		return
	}

	manifestJSONBytes, err := yaml.YAMLToJSON(ManifestYAMLBytes)
	if nil != err {
		// handle syntax errors specially
		errs = append(errs, fmt.Errorf("Error validating pkg. Details: %v", err.Error()))
		return
	}

	result, err := this.manifestSchema.Validate(
		gojsonschema.NewBytesLoader(manifestJSONBytes),
	)
	if nil != err {
		// handle syntax errors specially
		errs = append(errs, fmt.Errorf("Error validating pkg. Details: %v", err.Error()))
		return
	}
	for _, desc := range result.Errors() {
		errs = append(errs, fmt.Errorf("%s", desc))
	}
	return
}

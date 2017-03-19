package pkg

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/xeipuuv/gojsonschema"
	"path"
)

func (this pkg) Validate(
	pkgRef string,
) (errs []error) {
	ManifestYAMLBytes, err := this.fileSystem.GetBytesOfFile(
		path.Join(pkgRef, NameOfPackageManifestFile),
	)
	if nil != err {
		// handle syntax errors specially
		errs = append(errs, fmt.Errorf("Error validating pkg.\n Details: %v", err.Error()))
		return
	}

	manifestJSONBytes, err := yaml.YAMLToJSON(ManifestYAMLBytes)
	if nil != err {
		// handle syntax errors specially
		errs = append(errs, fmt.Errorf("Error validating pkg.\n Details: %v", err.Error()))
		return
	}

	result, err := this.manifestSchema.Validate(
		gojsonschema.NewBytesLoader(manifestJSONBytes),
	)
	if nil != err {
		// handle syntax errors specially
		errs = append(errs, fmt.Errorf("Error validating pkg.\n Details: %v", err.Error()))
		return
	}
	for _, desc := range result.Errors() {
		errs = append(errs, fmt.Errorf("%s\n", desc))
	}
	return
}

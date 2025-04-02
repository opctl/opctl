package platform

import (
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/str"
)

// Interpret OCI Image Platform
func Interpret(
	scope map[string]*model.Value,
	imagePlatformSpec *model.OCIImagePlatformSpec,
	scratchDir string,
) (*model.OCIImagePlatform, error) {
	ociImagePlatform := &model.OCIImagePlatform{}
	arch, err := str.Interpret(
		scope,
		imagePlatformSpec.Arch,
	)
	if err != nil {
		return nil, err
	}

	ociImagePlatform.Arch = arch.String

	return ociImagePlatform, nil
}

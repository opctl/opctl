package managepackages

import (
	"github.com/opspec-io/sdk-golang/pkg/model"
	"path"
)

func (this managePackages) SetPackageDescription(
	req model.SetPackageDescriptionReq,
) (err error) {

	pathToPackageManifestView := path.Join(req.PathToOp, NameOfPackageManifestFile)

	opBytes, err := this.fileSystem.GetBytesOfFile(
		pathToPackageManifestView,
	)
	if nil != err {
		return
	}

	packageManifestView := model.PackageManifestView{}
	err = this.yaml.To(
		opBytes,
		&packageManifestView,
	)
	if nil != err {
		return
	}

	packageManifestView.Description = req.Description

	opBytes, err = this.yaml.From(&packageManifestView)
	if nil != err {
		return
	}

	err = this.fileSystem.SaveFile(
		pathToPackageManifestView,
		opBytes,
	)

	return

}

package pkg

import (
	"github.com/opspec-io/sdk-golang/model"
	"path"
)

func (this pkg) SetDescription(
	req SetDescriptionReq,
) (err error) {

	pathToPackageManifestView := path.Join(req.Path, NameOfPackageManifestFile)

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

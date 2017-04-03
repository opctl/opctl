package pkg

import (
	"github.com/opspec-io/sdk-golang/model"
	"path"
)

func (this pkg) SetDescription(
	req SetDescriptionReq,
) (err error) {

	pathToPkgManifestView := path.Join(req.Path, NameOfPkgManifestFile)

	pkgManifestBytes, err := this.ioUtil.ReadFile(
		pathToPkgManifestView,
	)
	if nil != err {
		return
	}

	pkgManifestView := model.PackageManifestView{}
	err = this.yaml.To(
		pkgManifestBytes,
		&pkgManifestView,
	)
	if nil != err {
		return
	}

	pkgManifestView.Description = req.Description

	pkgManifestBytes, err = this.yaml.From(&pkgManifestView)
	if nil != err {
		return
	}

	err = this.ioUtil.WriteFile(
		pathToPkgManifestView,
		pkgManifestBytes,
		0777,
	)

	return

}

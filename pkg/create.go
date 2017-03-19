package pkg

import (
	"github.com/opspec-io/sdk-golang/model"
	"path"
)

func (this pkg) Create(
	req CreateReq,
) (err error) {

	err = this.fileSystem.AddDir(
		req.Path,
	)
	if nil != err {
		return
	}

	var packageManifestView = model.PackageManifestView{
		Description: req.Description,
		Name:        req.Name,
	}

	packageManifestViewBytes, err := this.yaml.From(&packageManifestView)
	if nil != err {
		return
	}

	err = this.fileSystem.SaveFile(
		path.Join(req.Path, NameOfPackageManifestFile),
		packageManifestViewBytes,
	)

	return

}

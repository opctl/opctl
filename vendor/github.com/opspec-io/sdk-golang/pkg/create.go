package pkg

import (
	"github.com/opspec-io/sdk-golang/model"
	"path"
)

func (this pkg) Create(
	req CreateReq,
) (err error) {

	err = this.fileSystem.MkdirAll(
		req.Path,
		0777,
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

	err = this.ioUtil.WriteFile(
		path.Join(req.Path, NameOfPkgManifestFile),
		packageManifestViewBytes,
		0777,
	)

	return

}

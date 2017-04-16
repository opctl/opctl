package pkg

import (
	"github.com/opspec-io/sdk-golang/model"
	"gopkg.in/yaml.v2"
	"path"
)

// Create creates an opspec package
func (this pkg) Create(
	req CreateReq,
) error {

	err := this.fs.MkdirAll(
		req.Path,
		0777,
	)
	if nil != err {
		return err
	}

	pkgManifest := model.PkgManifest{
		Description: req.Description,
		Name:        req.Name,
	}

	pkgManifestBytes, err := yaml.Marshal(&pkgManifest)
	if nil != err {
		return err
	}

	return this.ioUtil.WriteFile(
		path.Join(req.Path, OpDotYmlFileName),
		pkgManifestBytes,
		0777,
	)

}

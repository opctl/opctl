package pkg

import (
	"gopkg.in/yaml.v2"
	"path"
)

// SetDescription sets the description of a package
func (this pkg) SetDescription(
	req SetDescriptionReq,
) error {

	pkgManifest, err := this.manifestUnmarshaller.Unmarshal(req.Path)
	if nil != err {
		return err
	}

	pkgManifest.Description = req.Description

	pkgManifestBytes, err := yaml.Marshal(pkgManifest)
	if nil != err {
		return err
	}

	return this.ioUtil.WriteFile(
		path.Join(req.Path, OpDotYmlFileName),
		pkgManifestBytes,
		0777,
	)

}

package pkg

import (
	"gopkg.in/yaml.v2"
	"path"
)

// SetDescription sets the description of a package
func (this pkg) SetDescription(
	pkgPath,
	pkgDescription string,
) error {

	pkgManifest, err := this.manifestUnmarshaller.Unmarshal(pkgPath)
	if nil != err {
		return err
	}

	pkgManifest.Description = pkgDescription

	pkgManifestBytes, err := yaml.Marshal(pkgManifest)
	if nil != err {
		return err
	}

	return this.ioUtil.WriteFile(
		path.Join(pkgPath, OpDotYmlFileName),
		pkgManifestBytes,
		0777,
	)

}

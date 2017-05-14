package pkg

import (
	"gopkg.in/yaml.v2"
	"path/filepath"
)

// SetDescription sets the description of a package
func (this _Pkg) SetDescription(
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
		filepath.Join(pkgPath, OpDotYmlFileName),
		pkgManifestBytes,
		0777,
	)

}

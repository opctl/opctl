package pkg

import (
	"github.com/opspec-io/sdk-golang/model"
	"gopkg.in/yaml.v2"
	"path/filepath"
)

// Create creates an opspec package
func (this _Pkg) Create(
	path,
	pkgName,
	pkgDescription string,
) error {

	err := this.os.MkdirAll(
		path,
		0777,
	)
	if nil != err {
		return err
	}

	pkgManifest := model.PkgManifest{
		Description: pkgDescription,
		Name:        pkgName,
	}

	pkgManifestBytes, err := yaml.Marshal(&pkgManifest)
	if nil != err {
		return err
	}

	return this.ioUtil.WriteFile(
		filepath.Join(path, OpDotYmlFileName),
		pkgManifestBytes,
		0777,
	)

}

package pkg

//go:generate counterfeiter -o ./fakeList.go --fake-name fakeList ./ list

import (
	"github.com/opspec-io/sdk-golang/model"
	"path"
)

// List lists packages in a directory
func (this _Pkg) List(
	dirPath string,
) ([]*model.PkgManifest, error) {

	childFileInfos, err := this.ioUtil.ReadDir(dirPath)
	if nil != err {
		return nil, err
	}

	var pkgs []*model.PkgManifest
	for _, childFileInfo := range childFileInfos {
		pkgManifest, err := this.manifestUnmarshaller.Unmarshal(
			path.Join(dirPath, childFileInfo.Name()),
		)
		if nil == err {
			// ignore err'd pkgs
			pkgs = append(pkgs, pkgManifest)
		}
	}

	return pkgs, err
}

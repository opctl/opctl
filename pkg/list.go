package pkg

//go:generate counterfeiter -o ./fakeList.go --fake-name fakeList ./ list

import (
	"github.com/opspec-io/sdk-golang/model"
	"path/filepath"
)

// List recursively lists packages in dirPath
func (this _Pkg) List(
	dirPath string,
) ([]*model.PkgManifest, error) {

	childFileInfos, err := this.ioUtil.ReadDir(dirPath)
	if nil != err {
		return nil, err
	}

	var pkgs []*model.PkgManifest
	for _, childFileInfo := range childFileInfos {

		childPath := filepath.Join(dirPath, childFileInfo.Name())

		if childFileInfo.IsDir() {
			// recurse into child dirs
			childPkgs, err := this.List(childPath)
			if nil != err {
				return nil, err
			}
			pkgs = append(pkgs, childPkgs...)
		} else if childFileInfo.Name() == OpDotYmlFileName {
			manifestReader, err := this.os.Open(childPath)
			if nil == err {
				if manifest, err := this.manifest.Unmarshal(manifestReader); nil == err {
					manifestReader.Close()
					// ignore err'd pkgs
					pkgs = append(pkgs, manifest)
				}
			}
		}

	}

	return pkgs, err
}

package pkg

import (
	"github.com/opspec-io/sdk-golang/model"
	"path/filepath"
)

// ListOps recursively lists ops in dirPath
func (this _Pkg) ListOps(
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
			childPkgs, err := this.ListOps(childPath)
			if nil != err {
				return nil, err
			}
			pkgs = append(pkgs, childPkgs...)
		} else if childFileInfo.Name() == OpDotYmlFileName {

			manifestReader, err := this.os.Open(childPath)
			if nil != err {
				// ignore errors for now;
				continue
			}

			manifestBytes, err := this.ioUtil.ReadAll(manifestReader)
			manifestReader.Close()
			if nil != err {
				// ignore errors for now;
				continue
			}

			if manifest, err := this.manifest.Unmarshal(manifestBytes); nil == err {
				// ignore err'd pkgs
				pkgs = append(pkgs, manifest)
			}
		}

	}

	return pkgs, err
}

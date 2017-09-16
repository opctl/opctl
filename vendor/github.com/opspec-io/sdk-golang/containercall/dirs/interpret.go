package dirs

import (
	"fmt"
	"github.com/opspec-io/sdk-golang/model"
	"path/filepath"
	"strings"
)

func (d _Dirs) Interpret(
	pkgHandle model.PkgHandle,
	scope map[string]*model.Value,
	scgContainerCallDirs map[string]string,
	scratchDirPath string,
) (map[string]string, error) {
	dcgContainerCallDirs := map[string]string{}
	for scgContainerDirPath, dirExpression := range scgContainerCallDirs {

		if "" == dirExpression {
			// bound implicitly
			dirExpression = scgContainerDirPath
		}

		isBoundToPkg := strings.HasPrefix(dirExpression, "/")
		value, isBoundToScope := scope[dirExpression]

		switch {
		case isBoundToPkg:
			// bound to pkg dir
			dcgContainerCallDirs[scgContainerDirPath] = filepath.Join(scratchDirPath, dirExpression)

			// pkg dirs must be passed by value
			if err := d.dirCopier.OS(
				filepath.Join(pkgHandle.Ref(), dirExpression),
				dcgContainerCallDirs[scgContainerDirPath],
			); nil != err {
				return nil, err
			}
		case isBoundToScope:
			// bound to scope
			if nil == value || nil == value.Dir {
				return nil, fmt.Errorf(
					"unable to bind dir '%v' to '%v'; '%v' not a dir",
					scgContainerDirPath,
					dirExpression,
					dirExpression,
				)
			}

			if strings.HasPrefix(*value.Dir, d.rootFSPath) {
				// bound to rootFS dir
				dcgContainerCallDirs[scgContainerDirPath] = filepath.Join(scratchDirPath, scgContainerDirPath)

				// rootFS dirs must be passed by value
				if err := d.dirCopier.OS(
					*value.Dir,
					dcgContainerCallDirs[scgContainerDirPath],
				); nil != err {
					return nil, err
				}
			} else {
				// bound to non rootFS dir
				dcgContainerCallDirs[scgContainerDirPath] = *value.Dir
			}
		default:
			// unbound; create tree
			dcgContainerCallDirs[scgContainerDirPath] = filepath.Join(scratchDirPath, scgContainerDirPath)
			if err := d.os.MkdirAll(
				dcgContainerCallDirs[scgContainerDirPath],
				0700,
			); nil != err {
				return nil, err
			}
		}
	}
	return dcgContainerCallDirs, nil
}
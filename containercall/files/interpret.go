package files

import (
	"fmt"
	"github.com/opspec-io/sdk-golang/model"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func (f _Files) Interpret(
	pkgPath string,
	scope map[string]*model.Value,
	scgContainerCallFiles map[string]string,
	scratchDirPath string,
) (map[string]string, error) {
	dcgContainerCallFiles := map[string]string{}
	for scgContainerFilePath, scgContainerFileBind := range scgContainerCallFiles {

		if "" == scgContainerFileBind {
			// bound implicitly
			scgContainerFileBind = scgContainerFilePath
		}

		isBoundToPkg := strings.HasPrefix(scgContainerFileBind, "/")
		value, isBoundToScope := scope[scgContainerFileBind]

		switch {
		case isBoundToPkg:
			// bound to pkg file
			dcgContainerCallFiles[scgContainerFilePath] = filepath.Join(scratchDirPath, scgContainerFileBind)

			// pkg files must be passed by value
			if err := f.fileCopier.OS(
				filepath.Join(pkgPath, scgContainerFileBind),
				dcgContainerCallFiles[scgContainerFilePath],
			); nil != err {
				return nil, err
			}
		case isBoundToScope:
			// bound to scope
			if nil == value || (nil == value.File && nil == value.Number && nil == value.String) {
				return nil, fmt.Errorf(
					"Unable to bind file '%v' to '%v'. '%v' not a file, number, or string",
					scgContainerFilePath,
					scgContainerFileBind,
					scgContainerFileBind,
				)
			}

			if strings.HasPrefix(*value.File, f.rootFSPath) {
				// bound to rootFS file
				dcgContainerCallFiles[scgContainerFilePath] = filepath.Join(scratchDirPath, scgContainerFilePath)

				// rootFS files must be passed by value
				if err := f.fileCopier.OS(
					*value.File,
					dcgContainerCallFiles[scgContainerFilePath],
				); nil != err {
					return nil, err
				}
			} else {
				// bound to non rootFS file
				dcgContainerCallFiles[scgContainerFilePath] = *value.File
			}
		default:
			// unbound; create tree
			dcgContainerCallFiles[scgContainerFilePath] = filepath.Join(scratchDirPath, scgContainerFilePath)
			// create dir
			if err := f.os.MkdirAll(
				path.Dir(dcgContainerCallFiles[scgContainerFilePath]),
				0700,
			); nil != err {
				return nil, err
			}
			// create file
			var outputFile *os.File
			outputFile, err := f.os.Create(dcgContainerCallFiles[scgContainerFilePath])
			outputFile.Close()
			if nil != err {
				return nil, err
			}
		}
	}
	return dcgContainerCallFiles, nil
}

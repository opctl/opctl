package files

import (
	"fmt"
	"github.com/opspec-io/sdk-golang/model"
	"path/filepath"
	"strings"
)

func (f _Files) Interpret(
	pkgHandle model.PkgHandle,
	scope map[string]*model.Value,
	scgContainerCallFiles map[string]interface{},
	scratchDirPath string,
) (map[string]string, error) {
	dcgContainerCallFiles := map[string]string{}
fileLoop:
	for scgContainerFilePath, fileExpression := range scgContainerCallFiles {

		if "" == fileExpression {
			// bound implicitly
			fileExpression = fmt.Sprintf("$(%v)", scgContainerFilePath)
		}

		fileValue, err := f.expression.EvalToFile(
			scope,
			fileExpression,
			pkgHandle,
			scratchDirPath,
		)
		if nil != err {
			return nil, fmt.Errorf(
				"unable to bind %v to %v; error was %v",
				scgContainerFilePath,
				fileExpression,
				err,
			)
		}

		if !strings.HasPrefix(*fileValue.File, f.rootFSPath) {
			// bound to non rootFS file
			dcgContainerCallFiles[scgContainerFilePath] = *fileValue.File
			continue fileLoop
		}
		dcgContainerCallFiles[scgContainerFilePath] = filepath.Join(scratchDirPath, scgContainerFilePath)

		// create file
		if err := f.os.MkdirAll(
			filepath.Dir(dcgContainerCallFiles[scgContainerFilePath]),
			0777,
		); nil != err {
			return nil, fmt.Errorf(
				"unable to bind %v to %v; error was %v",
				scgContainerFilePath,
				fileExpression,
				err,
			)
		}

		err = f.fileCopier.OS(
			*fileValue.File,
			dcgContainerCallFiles[scgContainerFilePath],
		)
		if nil != err {
			return nil, fmt.Errorf(
				"unable to bind %v to %v; error was %v",
				scgContainerFilePath,
				fileExpression,
				err,
			)
		}

	}
	return dcgContainerCallFiles, nil
}

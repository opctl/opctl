package files

import (
	"fmt"
	"github.com/opspec-io/sdk-golang/model"
	"os"
	"path/filepath"
	"strings"
)

func (f _Files) Interpret(
	pkgHandle model.PkgHandle,
	scope map[string]*model.Value,
	scgContainerCallFiles map[string]string,
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

		file, err := f.os.Open(*fileValue.File)
		if nil != err {
			return nil, fmt.Errorf(
				"unable to bind %v to %v; error was %v",
				scgContainerFilePath,
				fileExpression,
				err,
			)
		}
		defer file.Close()

		fileInfo, err := f.os.Stat(*fileValue.File)
		if nil != err {
			return nil, fmt.Errorf(
				"unable to bind %v to %v; error was %v",
				scgContainerFilePath,
				fileExpression,
				err,
			)
		}

		outputFile, err := f.os.OpenFile(
			dcgContainerCallFiles[scgContainerFilePath],
			os.O_RDWR|os.O_CREATE,
			fileInfo.Mode(),
		)
		if nil != err {
			return nil, fmt.Errorf(
				"unable to bind %v to %v; error was %v",
				scgContainerFilePath,
				fileExpression,
				err,
			)
		}
		defer outputFile.Close()

		// copy content to file
		_, err = f.io.Copy(outputFile, file)
		file.Close()
		outputFile.Close()
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

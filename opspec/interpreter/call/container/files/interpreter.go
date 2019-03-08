package files

//go:generate counterfeiter -o ./fakeInterpreter.go --fake-name FakeInterpreter ./ Interpreter

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/golang-interfaces/ios"
	"github.com/golang-utils/filecopier"
	"github.com/opctl/sdk-golang/model"
	"github.com/opctl/sdk-golang/opspec/interpreter/file"
)

type Interpreter interface {
	Interpret(
		opHandle model.DataHandle,
		scope map[string]*model.Value,
		scgContainerCallFiles map[string]interface{},
		scratchDirPath string,
	) (map[string]string, error)
}

// NewInterpreter returns an initialized Interpreter instance
func NewInterpreter(
	dataDirPath string,
) Interpreter {
	return _interpreter{
		fileCopier:      filecopier.New(),
		fileInterpreter: file.NewInterpreter(),
		os:              ios.New(),
		dataDirPath:     dataDirPath,
	}
}

type _interpreter struct {
	fileCopier      filecopier.FileCopier
	fileInterpreter file.Interpreter
	os              ios.IOS
	dataDirPath     string
}

func (itp _interpreter) Interpret(
	opHandle model.DataHandle,
	scope map[string]*model.Value,
	scgContainerCallFiles map[string]interface{},
	scratchDirPath string,
) (map[string]string, error) {
	dcgContainerCallFiles := map[string]string{}
fileLoop:
	for scgContainerFilePath, fileExpression := range scgContainerCallFiles {

		if nil == fileExpression {
			// bound implicitly
			fileExpression = fmt.Sprintf("$(%v)", scgContainerFilePath)
		}

		fileValue, err := itp.fileInterpreter.Interpret(
			scope,
			fileExpression,
			opHandle,
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

		if !strings.HasPrefix(*fileValue.File, itp.dataDirPath) {
			// bound to non rootFS file
			dcgContainerCallFiles[scgContainerFilePath] = *fileValue.File
			continue fileLoop
		}
		dcgContainerCallFiles[scgContainerFilePath] = filepath.Join(scratchDirPath, scgContainerFilePath)

		// create file
		if err := itp.os.MkdirAll(
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

		err = itp.fileCopier.OS(
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

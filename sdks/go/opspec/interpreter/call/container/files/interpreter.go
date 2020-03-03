package files

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/golang-interfaces/ios"
	"github.com/golang-utils/filecopier"
	"github.com/opctl/opctl/sdks/go/data/coerce"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/file"
)

//counterfeiter:generate -o fakes/interpreter.go . Interpreter
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
		coerce:          coerce.New(),
		fileCopier:      filecopier.New(),
		fileInterpreter: file.NewInterpreter(),
		os:              ios.New(),
		dataDirPath:     dataDirPath,
	}
}

type _interpreter struct {
	coerce          coerce.Coerce
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
			// @TODO: return existence from fileInterpreter.Interpret (rather than treating all errors as due to non-existence) so we unambiguously know this is an assignment
			fileValue, err = itp.coerce.ToFile(&model.Value{String: new(string)}, scratchDirPath)
			if nil != err {
				return nil, fmt.Errorf(
					"unable to bind %v to %v; error was %v",
					scgContainerFilePath,
					fileExpression,
					err,
				)
			}
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

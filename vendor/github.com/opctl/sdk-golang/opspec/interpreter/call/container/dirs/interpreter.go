package dirs

//go:generate counterfeiter -o ./fakeInterpreter.go --fake-name FakeInterpreter ./ Interpreter

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/golang-interfaces/ios"
	"github.com/golang-utils/dircopier"
	"github.com/opctl/sdk-golang/model"
	"github.com/opctl/sdk-golang/opspec/interpreter/dir"
)

type Interpreter interface {
	Interpret(
		opHandle model.DataHandle,
		scope map[string]*model.Value,
		scgContainerCallFiles map[string]string,
		scratchDirPath string,
	) (map[string]string, error)
}

// NewInterpreter returns an initialized Interpreter instance
func NewInterpreter(
	dataDirPath string,
) Interpreter {
	return _interpreter{
		dirCopier:      dircopier.New(),
		dirInterpreter: dir.NewInterpreter(),
		os:             ios.New(),
		dataDirPath:    dataDirPath,
	}
}

type _interpreter struct {
	dirCopier      dircopier.DirCopier
	dirInterpreter dir.Interpreter
	os             ios.IOS
	dataDirPath    string
}

func (itp _interpreter) Interpret(
	opHandle model.DataHandle,
	scope map[string]*model.Value,
	scgContainerCallDirs map[string]string,
	scratchDirPath string,
) (map[string]string, error) {
	dcgContainerCallDirs := map[string]string{}
dirLoop:
	for scgContainerDirPath, dirExpression := range scgContainerCallDirs {

		if "" == dirExpression {
			// bound implicitly
			dirExpression = fmt.Sprintf("$(%v)", scgContainerDirPath)
		}

		dirValue, err := itp.dirInterpreter.Interpret(
			scope,
			dirExpression,
			opHandle,
		)
		if nil != err {
			return nil, fmt.Errorf(
				"unable to bind %v to %v; error was %v",
				scgContainerDirPath,
				dirExpression,
				err,
			)
		}

		if "" != *dirValue.Dir && !strings.HasPrefix(*dirValue.Dir, itp.dataDirPath) {
			// bound to non rootFS dir
			dcgContainerCallDirs[scgContainerDirPath] = *dirValue.Dir
			continue dirLoop
		}
		dcgContainerCallDirs[scgContainerDirPath] = filepath.Join(scratchDirPath, scgContainerDirPath)

		if "" == *dirValue.Dir {

			if err := itp.os.MkdirAll(
				dcgContainerCallDirs[scgContainerDirPath],
				0700,
			); nil != err {
				return nil, err
			}

		} else {

			if err := itp.dirCopier.OS(
				*dirValue.Dir,
				dcgContainerCallDirs[scgContainerDirPath],
			); nil != err {
				return nil, fmt.Errorf(
					"unable to bind %v to %v; error was %v",
					scgContainerDirPath,
					dirExpression,
					err,
				)
			}

		}
	}
	return dcgContainerCallDirs, nil
}

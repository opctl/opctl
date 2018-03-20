package dirs

//go:generate counterfeiter -o ./fakeInterpreter.go --fake-name FakeInterpreter ./ Interpreter

import (
	"fmt"
	"github.com/golang-interfaces/ios"
	"github.com/golang-utils/dircopier"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/op/interpreter/expression"
	"path/filepath"
	"strings"
)

type Interpreter interface {
	Interpret(
		opDirHandle model.DataHandle,
		scope map[string]*model.Value,
		scgContainerCallFiles map[string]string,
		scratchDirPath string,
	) (map[string]string, error)
}

// NewInterpreter returns an initialized Interpreter instance
func NewInterpreter(
	rootFSPath string,
) Interpreter {
	return _interpreter{
		dirCopier:  dircopier.New(),
		expression: expression.New(),
		os:         ios.New(),
		rootFSPath: rootFSPath,
	}
}

type _interpreter struct {
	dirCopier  dircopier.DirCopier
	expression expression.Expression
	os         ios.IOS
	rootFSPath string
}

func (itp _interpreter) Interpret(
	opDirHandle model.DataHandle,
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

		dirValue, err := itp.expression.EvalToDir(
			scope,
			dirExpression,
			opDirHandle,
		)
		if nil != err {
			return nil, fmt.Errorf(
				"unable to bind %v to %v; error was %v",
				scgContainerDirPath,
				dirExpression,
				err,
			)
		}

		if "" != *dirValue.Dir && !strings.HasPrefix(*dirValue.Dir, itp.rootFSPath) {
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

package dirs

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/golang-interfaces/ios"
	"github.com/golang-utils/dircopier"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/dir"
)

//counterfeiter:generate -o fakes/interpreter.go . Interpreter
type Interpreter interface {
	Interpret(
		scope map[string]*model.Value,
		callContainerSpecFiles map[string]string,
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
	scope map[string]*model.Value,
	callContainerSpecDirs map[string]string,
	scratchDirPath string,
) (map[string]string, error) {
	dcgContainerCallDirs := map[string]string{}
dirLoop:
	for callSpecContainerDirPath, dirExpression := range callContainerSpecDirs {

		if "" == dirExpression {
			// bound implicitly
			dirExpression = fmt.Sprintf("$(%v)", callSpecContainerDirPath)
		}

		dcgContainerCallDirs[callSpecContainerDirPath] = filepath.Join(scratchDirPath, callSpecContainerDirPath)
		dirValue, err := itp.dirInterpreter.Interpret(
			scope,
			dirExpression,
			scratchDirPath,
			true,
		)
		if nil != err {
			return nil, fmt.Errorf(
				"unable to bind %v to %v; error was %v",
				callSpecContainerDirPath,
				dirExpression,
				err,
			)
		}

		if "" != *dirValue.Dir && !strings.HasPrefix(*dirValue.Dir, itp.dataDirPath) {
			// bound to non rootFS dir
			dcgContainerCallDirs[callSpecContainerDirPath] = *dirValue.Dir
			continue dirLoop
		}

		if err := itp.dirCopier.OS(
			*dirValue.Dir,
			dcgContainerCallDirs[callSpecContainerDirPath],
		); nil != err {
			return nil, fmt.Errorf(
				"unable to bind %v to %v; error was %v",
				callSpecContainerDirPath,
				dirExpression,
				err,
			)
		}

	}
	return dcgContainerCallDirs, nil
}

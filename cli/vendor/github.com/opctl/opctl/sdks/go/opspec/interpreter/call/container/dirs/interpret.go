package dirs

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/golang-utils/dircopier"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/dir"
)

// Interpret container dirs
func Interpret(
	scope map[string]*model.Value,
	containerCallSpecDirs map[string]string,
	scratchDirPath string,
	dataDirPath string,
) (map[string]string, error) {
	containerCallDirs := map[string]string{}
dirLoop:
	for callSpecContainerDirPath, dirExpression := range containerCallSpecDirs {

		if "" == dirExpression {
			// bound implicitly
			dirExpression = fmt.Sprintf("$(%v)", callSpecContainerDirPath)
		}

		dirValue, err := dir.Interpret(
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

		if "" != *dirValue.Dir && !strings.HasPrefix(*dirValue.Dir, dataDirPath) {
			// bound to non rootFS dir
			containerCallDirs[callSpecContainerDirPath] = *dirValue.Dir
			continue dirLoop
		}
		
		containerCallDirs[callSpecContainerDirPath] = filepath.Join(scratchDirPath, callSpecContainerDirPath)

		dirCopier := dircopier.New()
		if err := dirCopier.OS(
			*dirValue.Dir,
			containerCallDirs[callSpecContainerDirPath],
		); nil != err {
			return nil, fmt.Errorf(
				"unable to bind %v to %v; error was %v",
				callSpecContainerDirPath,
				dirExpression,
				err,
			)
		}

	}
	return containerCallDirs, nil
}

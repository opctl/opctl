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
	containerCallSpecDirs map[string]interface{},
	scratchDirPath string,
	dataCachePath string,
) (map[string]string, error) {
	containerCallDirs := map[string]string{}
dirLoop:
	for callSpecContainerDirPath, dirExpression := range containerCallSpecDirs {

		if nil == dirExpression {
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

		if "" != *dirValue.Dir && !strings.HasPrefix(*dirValue.Dir, dataCachePath) {
			// bound to non dataCachePath
			containerCallDirs[callSpecContainerDirPath] = *dirValue.Dir
			continue dirLoop
		}

		// copy cached files to ensure can't be mutated
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

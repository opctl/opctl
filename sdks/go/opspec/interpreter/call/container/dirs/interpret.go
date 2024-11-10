package dirs

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/golang-utils/dircopier"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/dir"
)

// Interpret container dirs
func Interpret(
	scope map[string]*ipld.Node,
	containerCallSpecDirs map[string]interface{},
	scratchDirPath string,
	dataCachePath string,
) (model.StringMap, error) {
	containerCallDirs := map[string]string{}
dirLoop:
	for callSpecContainerDirPath, dirExpression := range containerCallSpecDirs {

		if dirExpression == nil {
			// bound implicitly
			dirExpression = opspec.NameToRef(callSpecContainerDirPath)
		}

		dirValue, err := dir.Interpret(
			scope,
			dirExpression,
			scratchDirPath,
			true,
		)
		if err != nil {
			return model.StringMap{}, fmt.Errorf("unable to bind directory %v to %v: %w", callSpecContainerDirPath, dirExpression, err)
		}

		if *dirValue.Dir != "" && !strings.HasPrefix(*dirValue.Dir, dataCachePath) {
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
		); err != nil {
			return model.StringMap{}, fmt.Errorf("unable to bind %v to %v: %w", callSpecContainerDirPath, dirExpression, err)
		}

	}
	return model.NewStringMap(
		containerCallDirs,
	), nil
}

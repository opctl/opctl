package files

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/opctl/opctl/sdks/go/internal/unsudo"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/file"
)

// Interpret container files
func Interpret(
	scope map[string]*model.Value,
	containerCallSpecFiles map[string]interface{},
	scratchDirPath string,
	dataCachePath string,
) (map[string]string, error) {
	containerCallFiles := map[string]string{}
fileLoop:
	for callSpecContainerFilePath, fileExpression := range containerCallSpecFiles {

		if fileExpression == nil {
			// bound implicitly
			fileExpression = opspec.NameToRef(callSpecContainerFilePath)
		}

		fileValue, err := file.Interpret(
			scope,
			fileExpression,
			scratchDirPath,
			true,
		)
		if err != nil {
			return nil, fmt.Errorf("unable to bind file %v to %v: %w", callSpecContainerFilePath, fileExpression, err)
		}

		if !strings.HasPrefix(*fileValue.File, dataCachePath) {
			// bound to non dataCachePath
			containerCallFiles[callSpecContainerFilePath] = *fileValue.File
			continue fileLoop
		}

		// copy cached files to ensure can't be mutated
		containerCallFiles[callSpecContainerFilePath] = filepath.Join(scratchDirPath, callSpecContainerFilePath)

		if err := unsudo.CloneFile(
			*fileValue.File,
			containerCallFiles[callSpecContainerFilePath],
		); err != nil {
			return nil, fmt.Errorf("unable to bind %v to %v: %w", callSpecContainerFilePath, fileExpression, err)
		}

	}
	return containerCallFiles, nil
}

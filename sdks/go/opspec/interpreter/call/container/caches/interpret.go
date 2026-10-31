package caches

import (
	"fmt"
	"path"
	"path/filepath"
	"regexp"

	"github.com/opctl/opctl/sdks/go/internal/unsudo"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/str"
)

var cacheNamePattern = regexp.MustCompile(`^[A-Za-z0-9._-]+$`)

// Interpret container cache dirs
func Interpret(
	scope map[string]*model.Value,
	containerCallSpecCaches map[string]interface{},
	dataDirPath string,
) (map[string]string, error) {
	containerCallDirs := map[string]string{}
	for specContainerPath, value := range containerCallSpecCaches {
		// clean so equivalent spellings (e.g. trailing slash) collide instead of failing late in the runtime
		containerPath := path.Clean(specContainerPath)
		if containerPath == "/" {
			return nil, fmt.Errorf("unable to bind cache dir %q: container path must not be /", specContainerPath)
		}
		if _, exists := containerCallDirs[containerPath]; exists {
			return nil, fmt.Errorf("unable to bind cache dir %q: container path is declared more than once", specContainerPath)
		}

		nameExpression, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("unable to bind cache dir %q: cache name must be a string expression", containerPath)
		}

		nameValue, err := str.Interpret(scope, nameExpression)
		if err != nil {
			return nil, fmt.Errorf("unable to bind cache dir %q: %w", containerPath, err)
		}
		name := *nameValue.String
		if !cacheNamePattern.MatchString(name) || name == "." || name == ".." {
			// truncate so a mistaken secret or large file reference isn't echoed in full
			shownName := name
			if len(shownName) > 64 {
				shownName = shownName[:64] + "..."
			}
			return nil, fmt.Errorf("unable to bind cache dir %q: invalid cache name %q: must match ^[A-Za-z0-9._-]+$ and must not be . or ..", containerPath, shownName)
		}

		hostPath := filepath.Join(dataDirPath, "caches", name)
		if err := unsudo.CreateDir(hostPath); err != nil {
			return nil, fmt.Errorf("unable to bind cache dir %q to %q: %w", containerPath, name, err)
		}
		containerCallDirs[containerPath] = hostPath
	}
	return containerCallDirs, nil
}

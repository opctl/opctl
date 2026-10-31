package caches

import (
	"fmt"
	"path"
	"path/filepath"
	"regexp"

	"github.com/opctl/opctl/sdks/go/internal/unsudo"
)

var cacheNamePattern = regexp.MustCompile(`^[A-Za-z0-9._-]+$`)

// Interpret container cache dirs
func Interpret(
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

		name, ok := value.(string)
		if !ok || name == "" {
			return nil, fmt.Errorf("unable to bind cache dir %q: cache name must be a non-empty string literal", containerPath)
		}

		if !cacheNamePattern.MatchString(name) || name == "." || name == ".." {
			return nil, fmt.Errorf("unable to bind cache dir %q: invalid cache name %q: must match ^[A-Za-z0-9._-]+$ and must not be . or ..", containerPath, name)
		}

		hostPath := filepath.Join(dataDirPath, "caches", name)
		if err := unsudo.CreateDir(hostPath); err != nil {
			return nil, fmt.Errorf("unable to bind cache dir %q to %q: %w", containerPath, name, err)
		}
		containerCallDirs[containerPath] = hostPath
	}
	return containerCallDirs, nil
}

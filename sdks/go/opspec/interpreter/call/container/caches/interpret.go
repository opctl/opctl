package caches

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/distribution/reference"
	"github.com/opctl/opctl/sdks/go/internal/unsudo"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/str"
)

var cacheNamePattern = regexp.MustCompile(`^[A-Za-z0-9._-]+$`)

// include hyphens so adjacent separators collapse to one
var cacheNameSeparatorPattern = regexp.MustCompile(`[^A-Za-z0-9._]+`)

// Interpret container cache dirs
func Interpret(
	scope map[string]*model.Value,
	containerCallSpecCaches map[string]interface{},
	imageRef *string,
	dataDirPath string,
	// noCacheRootPath, if set, replaces the cache root for a --no-cache run
	noCacheRootPath string,
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

		var name string
		if value == nil {
			if imageRef == nil {
				return nil, fmt.Errorf("unable to bind cache dir %q: image has no ref; give an explicit cache name for this container path", containerPath)
			}
			parsedRef, err := reference.ParseAnyReference(strings.ToLower(*imageRef))
			if err != nil {
				return nil, fmt.Errorf("unable to bind cache dir %q: unable to derive cache name from image ref: %w", containerPath, err)
			}
			namedRef, ok := parsedRef.(reference.Named)
			if !ok {
				return nil, fmt.Errorf("unable to bind cache dir %q: image ref has no repository; give an explicit cache name for this container path", containerPath)
			}
			name = strings.Trim(cacheNameSeparatorPattern.ReplaceAllString(reference.FamiliarName(namedRef)+containerPath, "-"), "-")
		} else {
			nameExpression, ok := value.(string)
			if !ok {
				return nil, fmt.Errorf("unable to bind cache dir %q: cache name must be a string expression", containerPath)
			}

			nameValue, err := str.Interpret(scope, nameExpression)
			if err != nil {
				return nil, fmt.Errorf("unable to bind cache dir %q: %w", containerPath, err)
			}
			name = *nameValue.String
		}
		if !cacheNamePattern.MatchString(name) || name == "." || name == ".." {
			// truncate so a mistaken secret or large file reference isn't echoed in full
			shownName := name
			if len(shownName) > 64 {
				shownName = shownName[:64] + "..."
			}
			return nil, fmt.Errorf("unable to bind cache dir %q: invalid cache name %q: must match ^[A-Za-z0-9._-]+$ and must not be . or ..", containerPath, shownName)
		}

		hostPath := filepath.Join(dataDirPath, "caches", name)
		if noCacheRootPath != "" {
			hostPath = filepath.Join(noCacheRootPath, name)
		}
		if err := unsudo.CreateDir(hostPath); err != nil {
			return nil, fmt.Errorf("unable to bind cache dir %q to %q: %w", containerPath, name, err)
		}
		info, err := os.Lstat(hostPath)
		if err != nil {
			return nil, fmt.Errorf("unable to bind cache dir %q to %q: %w", containerPath, name, err)
		}
		if !info.IsDir() {
			return nil, fmt.Errorf("unable to bind cache dir %q to %q: mount source %q must be a real directory", containerPath, name, hostPath)
		}
		// Best effort: an unprivileged node cannot chmod another user's cache dir.
		if info.Mode().Perm() != 0777 {
			if err := os.Chmod(hostPath, info.Mode()|0777); err != nil && !errors.Is(err, fs.ErrPermission) {
				return nil, fmt.Errorf("unable to bind cache dir %q to %q: %w", containerPath, name, err)
			}
		}
		containerCallDirs[containerPath] = hostPath
	}
	return containerCallDirs, nil
}

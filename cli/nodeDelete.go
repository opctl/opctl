package main

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/opctl/opctl/cli/internal/nodeprovider/local"
)

func removeAllExcept(rootDir string, keepPaths ...string) error {
	if len(keepPaths) == 0 {
		return os.RemoveAll(rootDir)
	}

	keepSet := make(map[string]struct{})
	for _, kp := range keepPaths {
		kp = filepath.Clean(kp)
		keepSet[kp] = struct{}{}
	}

	directories := []string{}
	pathsKept := []string{}
	err := filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		relPath, err := filepath.Rel(rootDir, path)
		if err != nil {
			return err
		}

		if relPath != "." {
			if _, keep := keepSet[relPath]; keep {
				pathsKept = append(pathsKept, path)
				return nil
			}
		}

		if d.IsDir() {
			directories = append(directories, path)
		} else {
			err = os.Remove(path)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	for _, dir := range directories {
		if len(pathsKept) == 0 {
			os.RemoveAll(dir)
			continue
		}
		for _, path := range pathsKept {
			if isSub, _ := subElement(dir, path); !isSub {
				os.RemoveAll(dir)
				break
			}
		}
	}

	return nil
}

func subElement(parent, sub string) (bool, error) {
	up := ".." + string(os.PathSeparator)

	// path-comparisons using filepath.Abs don't work reliably according to docs (no unique representation).
	rel, err := filepath.Rel(parent, sub)
	if err != nil {
		return false, err
	}
	if !strings.HasPrefix(rel, up) && rel != ".." {
		return true, nil
	}
	return false, nil
}

// nodeDelete implements the node delete command
func nodeDelete(
	ctx context.Context,
	nodeConfig local.NodeConfig,
) error {
	containerRT, err := getContainerRuntime(ctx, nodeConfig)
	if err != nil {
		return err
	}

	np, err := local.New(nodeConfig)
	if err != nil {
		return err
	}

	err = containerRT.Delete(ctx)
	if err != nil {
		return err
	}

	if err := np.KillNodeIfExists(); err != nil {
		return err
	}

	// In macOS XDG_CONFIG_HOME points to the same directory as XDG_DATA_HOME,
	// we can remove the whole directory except the telemetry/config.yaml file
	return removeAllExcept(nodeConfig.DataDir, "telemetry/config.yaml")
}

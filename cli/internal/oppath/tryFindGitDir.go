package oppath

import (
	"fmt"
	"os"
	"path/filepath"
)

func tryFindGitDir(
	currentPath string,
) (string, error) {
	currentPath, err := filepath.Abs(currentPath)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %w", err)
	}

	for {
		gitPath := filepath.Join(currentPath, ".git")
		if info, err := os.Stat(gitPath); err == nil && info.IsDir() {
			return gitPath, nil
		}

		parentPath := filepath.Dir(currentPath)
		if parentPath == currentPath {
			break
		}
		currentPath = parentPath
	}

	return "", nil
}

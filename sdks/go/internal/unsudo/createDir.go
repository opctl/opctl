package unsudo

import (
	"os"
	"path/filepath"
	"strings"
)

// CreateDir as the user & group who ran sudo
func CreateDir(
	fsPath string,
) error {
	// Split the path into individual components
	parts := strings.Split(fsPath, string(os.PathSeparator))

	// Iterate over the path components, creating directories as needed
	currentPath := string(os.PathSeparator)
	for i := 0; i < len(parts); i++ {
		currentPath = filepath.Join(currentPath, parts[i])

		if err := os.Mkdir(currentPath, 0700); err != nil {
			if os.IsExist(err) {
				// if existing nothing to do
				continue
			}

			return err
		}

		if err := os.Chown(currentPath, getSudoUID(), getSudoGID()); err != nil {
			return err
		}
	}

	return nil
}

package unsudo

import (
	"io/fs"
	"path/filepath"
)

// CloneDir as the user & group who ran sudo
func CloneDir(
	srcPath,
	dstPath string,
) error {
	return filepath.WalkDir(
		srcPath,
		func(entrySrcPath string, entry fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			entryRelPath, err := filepath.Rel(srcPath, entrySrcPath)
			if err != nil {
				return err
			}

			entryDstPath := filepath.Join(dstPath, entryRelPath)

			if entry.IsDir() {
				return CreateDir(entryDstPath)
			}

			return CloneFile(entrySrcPath, entryDstPath)
		},
	)
}

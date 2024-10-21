package dns

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func deleteAllResolvers() error {
	return filepath.WalkDir(
		etcResolverPath,
		func(path string, d fs.DirEntry, err error) error {
			if err != nil && !os.IsNotExist(err) {
				if os.IsNotExist(err) {
					return nil
				}

				return err
			}

			if strings.HasPrefix(
				d.Name(),
				"opctl_",
			) {
				return os.Remove(path)
			}

			return nil
		},
	)
}

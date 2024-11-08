package resolvercfg

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Delete modifications to the current system
func Delete(
	ctx context.Context,
) error {
	if err := filepath.WalkDir(
		resolverDir,
		func(path string, d fs.DirEntry, err error) error {
			if err != nil && !os.IsNotExist(err) {
				if os.IsNotExist(err) {
					return nil
				}

				return err
			}

			if strings.HasPrefix(
				d.Name(),
				resolverPrefix,
			) {
				return os.Remove(path)
			}

			return nil
		},
	); err != nil {
		return err
	}

	return clearCaches(ctx)
}

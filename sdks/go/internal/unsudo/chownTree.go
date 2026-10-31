package unsudo

import (
	"errors"
	"io/fs"
	"os"
)

// ChownTree to the user & group who ran sudo, without following symlinks.
// Resolution is confined to fsPath so a concurrently swapped symlink can't redirect the walk.
func ChownTree(fsPath string) error {
	root, err := os.OpenRoot(fsPath)
	if err != nil {
		return err
	}
	defer root.Close()

	// shortcut: re-chowns entries that are already owned correctly; skip by stat if large caches make this slow.
	uid, gid := getSudoUID(), getSudoGID()
	var errs []error
	walkErr := fs.WalkDir(
		root.FS(),
		".",
		func(entryPath string, _ fs.DirEntry, err error) error {
			if err == nil {
				err = root.Lchown(entryPath, uid, gid)
			}
			// entries may vanish under calls sharing the cache dir
			if err != nil && !errors.Is(err, fs.ErrNotExist) {
				errs = append(errs, err)
			}
			return nil
		},
	)
	return errors.Join(append(errs, walkErr)...)
}

package unsudo

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
)

// RestoreCacheTree adds a+rwX and, only under sudo, restores ownership.
// Resolution is confined to fsPath so a concurrently swapped symlink can't redirect the walk.
func RestoreCacheTree(fsPath string) error {
	info, err := os.Lstat(fsPath)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("cache path %q must be a real directory", fsPath)
	}
	root, err := os.OpenRoot(fsPath)
	if err != nil {
		return err
	}
	defer root.Close()
	openedInfo, err := root.Stat(".")
	if err != nil {
		return err
	}
	if !os.SameFile(info, openedInfo) {
		return fmt.Errorf("cache path %q changed while opening its root", fsPath)
	}

	uid, gid := getSudoUID(), getSudoGID()
	chown := IsSudo()
	var errs []error
	walkErr := fs.WalkDir(
		root.FS(),
		".",
		func(entryPath string, entry fs.DirEntry, err error) error {
			if err == nil {
				var info fs.FileInfo
				info, err = entry.Info()
				// a hard-linked file may also live outside the cache dir; never re-own or widen it
				// shortcut: link check is by path, so a hard link swapped in after Lstat can still be re-owned/widened; use fd-based fchown/fchmod if hosts lack protected_hardlinks
				if err == nil && !isHardLinked(info) {
					if chown && !isOwnedBy(info, uid, gid) {
						if chownErr := root.Lchown(entryPath, uid, gid); chownErr != nil && !errors.Is(chownErr, fs.ErrNotExist) {
							errs = append(errs, chownErr)
						}
					}
					err = makeCacheEntryWritable(root, entryPath, info)
				}
			}
			// entries may vanish under calls sharing the cache dir
			if err == nil || errors.Is(err, fs.ErrNotExist) {
				return nil
			}
			// an unprivileged node can't read or adjust what other users (e.g. root containers) wrote; skip it silently
			if errors.Is(err, fs.ErrPermission) && os.Geteuid() != 0 {
				if entry != nil && entry.IsDir() {
					return fs.SkipDir
				}
				return nil
			}
			errs = append(errs, err)
			return nil
		},
	)
	return errors.Join(append(errs, walkErr)...)
}

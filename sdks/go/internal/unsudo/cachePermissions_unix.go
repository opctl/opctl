//go:build unix

package unsudo

import (
	"io/fs"
	"os"
	"syscall"
)

func makeCacheEntryWritable(root *os.Root, entryPath string, info fs.FileInfo) error {
	mode := info.Mode()
	if !mode.IsDir() && !mode.IsRegular() {
		return nil
	}
	// An unprivileged node cannot chmod another user's entries; skip them without warnings.
	if stat, ok := info.Sys().(*syscall.Stat_t); ok && os.Geteuid() != 0 && stat.Uid != uint32(os.Geteuid()) {
		return nil
	}
	target := mode | 0666
	if mode.IsDir() || mode.Perm()&0111 != 0 {
		target |= 0111
	}
	if mode.IsRegular() {
		// World-writable executables must not retain privilege-granting bits.
		target &^= os.ModeSetuid | os.ModeSetgid
	}
	if mode == target {
		return nil
	}

	// Root.Chmod confines symlink resolution and handles kernels without fchmodat2.
	// shortcut: target comes from the earlier Lstat, so an entry swapped mid-walk by a container sharing the cache can get a stale mode until the next call; recheck after chmod if that matters.
	if err := root.Chmod(entryPath, target); err != nil {
		if current, statErr := root.Lstat(entryPath); statErr == nil && current.Mode()&os.ModeSymlink != 0 {
			return nil
		}
		return err
	}
	return nil
}

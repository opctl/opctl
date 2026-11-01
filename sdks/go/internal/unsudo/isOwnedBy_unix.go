//go:build unix

package unsudo

import (
	"io/fs"
	"syscall"
)

func isOwnedBy(info fs.FileInfo, uid, gid int) bool {
	stat, ok := info.Sys().(*syscall.Stat_t)
	return ok && int(stat.Uid) == uid && int(stat.Gid) == gid
}

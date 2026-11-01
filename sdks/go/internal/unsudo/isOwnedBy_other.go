//go:build !unix

package unsudo

import "io/fs"

func isOwnedBy(info fs.FileInfo, uid, gid int) bool {
	return false
}

func isHardLinked(info fs.FileInfo) bool {
	return false
}

//go:build !unix

package unsudo

import (
	"io/fs"
	"os"
)

func makeCacheEntryWritable(root *os.Root, entryPath string, info fs.FileInfo) error {
	return nil
}

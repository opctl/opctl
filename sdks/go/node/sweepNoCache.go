package node

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Called with the event DB lock held, before any new runs can create nocache dirs.
func sweepNoCache(dataDirPath string) {
	sweepDir, err := os.MkdirTemp(dataDirPath, "nocache-sweep-")
	if err == nil {
		err = os.Rename(filepath.Join(dataDirPath, "nocache"), filepath.Join(sweepDir, "nocache"))
	}
	if err != nil && !os.IsNotExist(err) {
		// e.g. nocache is a mount point; remove in place so leftovers can't build up
		if err := os.RemoveAll(filepath.Join(dataDirPath, "nocache")); err != nil {
			fmt.Printf("unable to remove leftover --no-cache dirs: %v\n", err)
		}
	}

	entries, err := os.ReadDir(dataDirPath)
	if err != nil {
		fmt.Printf("unable to list leftover --no-cache dirs: %v\n", err)
		return
	}
	var sweepPaths []string
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), "nocache-sweep-") {
			sweepPaths = append(sweepPaths, filepath.Join(dataDirPath, entry.Name()))
		}
	}
	go func() {
		for _, sweepPath := range sweepPaths {
			if err := os.RemoveAll(sweepPath); err != nil {
				fmt.Printf("unable to remove leftover --no-cache dirs: %v\n", err)
			}
		}
	}()
}

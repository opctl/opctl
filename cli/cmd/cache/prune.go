package cache

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/opctl/opctl/cli/internal/euid0"
	"github.com/opctl/opctl/cli/internal/nodeprovider/local"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/container/caches"
	"github.com/spf13/cobra"
)

func newPruneCmd(nodeConfig *local.NodeConfig) *cobra.Command {
	return &cobra.Command{
		Use:   "prune [NAME...]",
		Short: "Delete cache dirs",
		Long:  "Delete the named cache dirs, or empty the cache root when no names are given. Only prune when no call is using the affected dirs. Missing names are ignored. Scratch dirs from --no-cache runs are left untouched.",
		Args: func(cmd *cobra.Command, args []string) error {
			return validateNames(args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := euid0.Ensure(); err != nil {
				return err
			}
			return prune(nodeConfig.DataDir, args)
		},
	}
}

func validateNames(names []string) error {
	for _, name := range names {
		if err := caches.ValidateName(name); err != nil {
			return err
		}
	}
	return nil
}

func prune(dataDir string, names []string) error {
	// Validate the whole request before removing anything, including when called without Cobra.
	if err := validateNames(names); err != nil {
		return err
	}

	cacheRoot := filepath.Join(dataDir, "caches")
	info, err := os.Lstat(cacheRoot)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("cache root %q must be a real directory", cacheRoot)
	}

	// remove through an os.Root so a cache root swapped for a symlink after the check can't redirect deletes
	root, err := os.OpenRoot(cacheRoot)
	if err != nil {
		return err
	}
	defer root.Close()
	if rootInfo, err := root.Stat("."); err != nil || !os.SameFile(info, rootInfo) {
		return fmt.Errorf("cache root %q changed while pruning", cacheRoot)
	}

	if len(names) == 0 {
		entries, err := fs.ReadDir(root.FS(), ".")
		if err != nil {
			return err
		}
		for _, entry := range entries {
			names = append(names, entry.Name())
		}
	}
	for _, name := range names {
		if err := root.RemoveAll(name); err != nil {
			return err
		}
	}
	return nil
}

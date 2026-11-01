package op

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/opctl/opctl/cli/internal/opstubs"
	"github.com/opctl/opctl/sdks/go/data"
	datafs "github.com/opctl/opctl/sdks/go/data/fs"
	"github.com/opctl/opctl/sdks/go/opspec/goop"
	"github.com/spf13/cobra"
)

func newGenCmd() *cobra.Command {
	return &cobra.Command{Use: "gen [OP_REF...]", Short: "Refresh op stubs in opctlgen/ for local Go modules", RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		roots := map[string]bool{}
		var failures []error
		add := func(dir string) {
			root, err := goop.ModuleRoot(dir)
			if err != nil {
				failures = append(failures, err)
			} else {
				roots[root] = true
			}
		}
		if len(args) == 0 {
			err := filepath.WalkDir(filepath.Join(cwd, ".opspec"), func(path string, entry fs.DirEntry, err error) error {
				if os.IsNotExist(err) {
					return nil
				}
				if err != nil {
					return err
				}
				if entry.IsDir() && entry.Name() == "opctlgen" {
					return filepath.SkipDir
				}
				if !entry.IsDir() && entry.Name() == "op.go" {
					add(filepath.Dir(path))
				}
				return nil
			})
			if err != nil {
				failures = append(failures, err)
			}
		} else {
			for _, ref := range args {
				if strings.Contains(ref, "#") {
					failures = append(failures, fmt.Errorf("op gen needs local op refs, got %q", ref))
					continue
				}
				dir, err := data.Resolve(cmd.Context(), ref, datafs.New(filepath.Join(cwd, ".opspec"), cwd))
				if err != nil {
					failures = append(failures, err)
					continue
				}
				add(dir.Ref())
			}
		}
		for root := range roots {
			if err := opstubs.RefreshModule(cmd.Context(), root, node); err != nil {
				failures = append(failures, err)
			}
		}
		return errors.Join(failures...)
	}}
}

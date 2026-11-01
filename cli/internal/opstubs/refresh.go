// Package opstubs refreshes editor stubs through the node's describe API.
package opstubs

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/goop"
)

type Describer interface {
	DescribeOp(context.Context, model.DescribeOpReq) (*model.OpSpec, error)
}

// RefreshModule refreshes every declaration, writing successes even on failure.
func RefreshModule(ctx context.Context, root string, node Describer) error {
	uses, err := goop.ScanModuleStubs(root)
	var failures []error
	if err != nil {
		failures = append(failures, err)
	}
	files := map[string][]byte{}
	for _, use := range uses {
		files[use.Name] = nil
		spec, err := node.DescribeOp(ctx, model.DescribeOpReq{Op: model.StartOpReqOp{Ref: goop.StubModuleRef(root, use)}})
		if err == nil {
			files[use.Name], err = goop.RenderStub(use, spec)
		}
		if err != nil {
			failures = append(failures, fmt.Errorf("%s: op stub %s (%s): %w", use.Site, use.Name, use.Ref, err))
		}
	}
	// A scanner failure without a usable declaration set must not delete stubs.
	if err == nil || uses != nil {
		if err := goop.WriteStubs(root, files); err != nil {
			failures = append(failures, err)
		}
	}
	return errors.Join(failures...)
}

// WarnRefresh refreshes only local Go source outside the node's data directory.
func WarnRefresh(ctx context.Context, ref string, node Describer, dataDir string, warnings io.Writer) {
	if strings.Contains(ref, "#") || !filepath.IsAbs(ref) {
		return
	}
	if _, err := os.Stat(filepath.Join(ref, "op.go")); err != nil {
		return
	}
	if dataDir != "" {
		path, err := filepath.EvalSymlinks(ref)
		if err != nil {
			return
		}
		bound, err := filepath.EvalSymlinks(dataDir)
		if err == nil {
			rel, err := filepath.Rel(bound, path)
			if err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
				return
			}
		}
	}
	root, err := goop.ModuleRoot(ref)
	if err == nil {
		err = RefreshModule(ctx, root, node)
	}
	if err != nil {
		message := strings.Join(strings.Fields(err.Error()), " ")
		fmt.Fprintf(warnings, "warning: opctlgen/ not refreshed: %s; run \"opctl op gen\" for details\n", message)
	}
}

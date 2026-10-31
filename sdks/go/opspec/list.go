package opspec

import (
	"context"
	"fmt"
	"io"
	"path/filepath"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/opfile"
)

// List ops recursively within a directory, describing Go ops through the node.
func List(
	ctx context.Context,
	dirHandle model.DataHandle,
	describe func(context.Context, string) (*model.OpSpec, error),
) (map[string]*model.OpSpec, error) {

	contents, err := dirHandle.ListDescendants(ctx)
	if err != nil {
		return nil, err
	}

	paths := map[string]bool{}
	for _, content := range contents {
		paths[filepath.Clean(content.Path)] = true
	}
	for path := range paths {
		if filepath.Base(path) == opfile.FileName && paths[filepath.Join(filepath.Dir(path), "op.go")] {
			return nil, opfile.ErrBothOpFiles
		}
	}

	opsByPath := map[string]*model.OpSpec{}
	for _, content := range contents {
		if filepath.Base(content.Path) == "op.go" {
			if describe == nil {
				return nil, fmt.Errorf("listing Go ops needs a node describe function")
			}
			path := filepath.Dir(content.Path)
			spec, err := describe(ctx, path)
			if err != nil {
				return nil, fmt.Errorf("error describing %s: %w", path, err)
			}
			opsByPath[path] = spec
		}
		if filepath.Base(content.Path) == opfile.FileName {

			opFileReader, err := dirHandle.GetContent(ctx, content.Path)
			if err != nil {
				return nil, fmt.Errorf("error opening %s%s: %w", dirHandle.Ref(), content.Path, err)
			}

			opFileBytes, err := io.ReadAll(opFileReader)
			opFileReader.Close()
			if err != nil {
				return nil, fmt.Errorf("error reading %s%s: %w", dirHandle.Ref(), content.Path, err)
			}

			opFile, err := opfile.Unmarshal(
				opFileBytes,
			)
			if err != nil {
				return nil, fmt.Errorf("error unmarshalling %s%s: %w", dirHandle.Ref(), content.Path, err)
			}

			opsByPath[filepath.Dir(content.Path)] = opFile
		}

	}

	return opsByPath, nil
}

package opspec

import (
	"context"
	"fmt"
	"io"
	"path/filepath"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/opfile"
)

// List ops recursively within a directory, returning discovered op files by path.
func List(
	ctx context.Context,
	dirHandle model.DataHandle,
) (map[string]*model.OpSpec, error) {

	contents, err := dirHandle.ListDescendants(ctx)
	if err != nil {
		return nil, err
	}

	opsByPath := map[string]*model.OpSpec{}
	for _, content := range contents {
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

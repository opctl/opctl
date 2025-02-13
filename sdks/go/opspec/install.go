package opspec

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/opctl/opctl/sdks/go/model"
	"golang.org/x/sync/errgroup"
)

// Install an op at path
func Install(
	ctx context.Context,
	path string,
	handle model.DataHandle,
) error {
	contentsList, err := handle.ListDescendants(ctx)
	if err != nil {
		return err
	}

	eg, ctx := errgroup.WithContext(ctx)
	eg.SetLimit(250)

	for _, content := range contentsList {
		eg.Go(func() error {

			dstPath := filepath.Join(path, content.Path)

			if _, statErr := os.Stat(dstPath); statErr == nil {
				// don't overwrite existing content
				return nil
			} else if !os.IsNotExist(statErr) {
				return statErr
			}

			if content.Mode.IsDir() {
				// ensure content path exists
				return os.MkdirAll(dstPath, content.Mode)
			}

			// ensure content dir exists
			if err := os.MkdirAll(filepath.Dir(dstPath), 0777); err != nil {
				return err
			}

			dst, err := os.Create(dstPath)
			if err != nil {
				return err
			}
			defer dst.Close()

			if err := os.Chmod(dstPath, content.Mode); err != nil {
				return err
			}

			src, err := handle.GetContent(ctx, content.Path)
			if err != nil {
				return err
			}
			defer src.Close()

			_, err = io.Copy(dst, src)
			return err
		})
	}

	return eg.Wait()
}

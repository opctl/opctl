package opspec

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/opctl/opctl/sdks/go/model"
)

// Install an op at path
func Install(
	ctx context.Context,
	path string,
	handle model.DataHandle,
) error {
	contentsList, err := handle.ListDescendants(ctx)
	if nil != err {
		return err
	}

	for _, content := range contentsList {
		dstPath := filepath.Join(path, content.Path)

		if _, statErr := os.Stat(dstPath); nil == statErr {
			// don't overwrite existing content
			continue
		} else if !os.IsNotExist(statErr) {
			return statErr
		}

		if content.Mode.IsDir() {
			// ensure content path exists
			if err = os.MkdirAll(dstPath, content.Mode); err != nil {
				return err
			}
		} else {
			// ensure content dir exists
			if err = os.MkdirAll(filepath.Dir(dstPath), 0777); err != nil {
				return err
			}

			dst, err := os.Create(dstPath)
			if nil != err {
				return err
			}
			defer dst.Close()

			if err = os.Chmod(dstPath, content.Mode); err != nil {
				return err
			}

			src, err := handle.GetContent(ctx, content.Path)
			if nil != err {
				return err
			}
			defer src.Close()

			if _, err = io.Copy(dst, src); err != nil {
				return err
			}
		}
	}

	return nil
}

package pkg

import (
	"context"
	"github.com/opspec-io/sdk-golang/model"
	"io"
	"path/filepath"
)

// Install installs an opspec pkg at path
func (this _Pkg) Install(
	ctx context.Context,
	path string,
	handle model.DataHandle,
) error {
	contentsList, err := handle.ListContents(ctx)
	if nil != err {
		return err
	}

	for _, content := range contentsList {

		dstPath := filepath.Join(path, content.Path)

		if content.Mode.IsDir() {
			// ensure content path exists
			err = this.os.MkdirAll(
				dstPath,
				content.Mode,
			)
			if nil != err {
				return err
			}
		} else {
			// ensure content dir exists
			err = this.os.MkdirAll(
				filepath.Dir(dstPath),
				0777,
			)
			if nil != err {
				return err
			}

			dst, err := this.os.Create(dstPath)
			if nil != err {
				return err
			}

			err = this.os.Chmod(dstPath, content.Mode)
			if nil != err {
				return err
			}

			src, err := handle.GetContent(ctx, content.Path)
			if nil != err {
				return err
			}

			_, err = io.Copy(dst, src)
			src.Close()
			dst.Close()
		}
	}

	return err

}

package pkg

import (
	"context"
	"io"
	"path/filepath"
)

// Install installs an opspec pkg at path
func (this _Pkg) Install(
	ctx context.Context,
	path string,
	handle Handle,
) error {

	err := this.os.MkdirAll(
		path,
		0777,
	)
	if nil != err {
		return err
	}

	contentsList, err := handle.ListContents(ctx)
	if nil != err {
		return err
	}

	for _, content := range contentsList {
		src, err := handle.GetContent(ctx, content.Path)
		if nil != err {
			return err
		}

		dstPath := filepath.Join(path, handle.Ref(), content.Path)

		// ensure content path exists
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

		_, err = io.Copy(dst, src)
		dst.Close()
	}

	return err

}

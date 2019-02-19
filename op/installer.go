package op

//go:generate counterfeiter -o ./fakeInstaller.go --fake-name FakeInstaller ./ Installer

import (
	"context"
	"github.com/golang-interfaces/ios"
	"github.com/opctl/sdk-golang/model"
	"io"
	"path/filepath"
)

type Installer interface {
	// Install installs an op; path will be created if it doesn't exist
	Install(
		ctx context.Context,
		path string,
		opHandle model.DataHandle,
	) error
}

// NewInstaller returns an initialized Installer instance
func NewInstaller() Installer {
	return _installer{
		os: ios.New(),
	}
}

type _installer struct {
	os ios.IOS
}

// Install installs an opspec pkg at path
func (inst _installer) Install(
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

		if content.Mode.IsDir() {
			// ensure content path exists
			err = inst.os.MkdirAll(
				dstPath,
				content.Mode,
			)
			if nil != err {
				return err
			}
		} else {
			// ensure content dir exists
			err = inst.os.MkdirAll(
				filepath.Dir(dstPath),
				0777,
			)
			if nil != err {
				return err
			}

			dst, err := inst.os.Create(dstPath)
			if nil != err {
				return err
			}

			err = inst.os.Chmod(dstPath, content.Mode)
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

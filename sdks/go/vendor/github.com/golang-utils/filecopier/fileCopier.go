package filecopier

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ FileCopier

import (
	"github.com/golang-interfaces/ios"
	"io"
)

type FileCopier interface {
	// copies an OS file from srcPath to dstPath. Creates or overwrites the destination as needed.
	OS(srcPath string, dstPath string) (err error)
}

func New() FileCopier {
	return fileCopier{
		os: ios.New(),
	}
}

type fileCopier struct {
	os ios.IOS
}

func (fc fileCopier) OS(srcPath string, dstPath string) (err error) {
	srcFile, err := fc.os.Open(srcPath)
	if nil != err {
		return
	}
	defer srcFile.Close()

	srcFileInfo, err := fc.os.Stat(srcPath)
	if nil != err {
		return
	}

	// copy content
	writer, err := fc.os.Create(dstPath)
	if nil != err {
		return
	}
	defer writer.Close()

	// copy mode
	err = fc.os.Chmod(dstPath, srcFileInfo.Mode())
	if nil != err {
		return
	}

	_, err = io.Copy(writer, srcFile)

	err = writer.Sync()

	return
}

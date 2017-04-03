package dircopier

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ DirCopier

import (
	"fmt"
	"github.com/virtual-go/filecopier"
	"github.com/virtual-go/fs"
	"github.com/virtual-go/fs/osfs"
	"github.com/virtual-go/vioutil"
	"path"
)

type DirCopier interface {
	// copies a fs dir from srcPath to dstPath. Creates or overwrites the destination as needed.
	Fs(srcPath string, dstPath string) (err error)
}

func New() DirCopier {
	_fs := osfs.New()
	return dirCopier{
		fs:         _fs,
		ioutil:     vioutil.New(_fs),
		fileCopier: filecopier.New(),
	}
}

type dirCopier struct {
	fs         fs.FS
	ioutil     vioutil.VIOUtil
	fileCopier filecopier.FileCopier
}

func (this dirCopier) Fs(srcPath string, dstPath string) (err error) {
	// get properties of srcPath
	fi, err := this.fs.Stat(srcPath)
	if nil != err {
		return
	}

	if !fi.IsDir() {
		err = fmt.Errorf("%v is not a dir", srcPath)
		return
	}

	// create dstPath
	err = this.fs.MkdirAll(dstPath, fi.Mode())
	if nil != err {
		return
	}

	entries, err := this.ioutil.ReadDir(srcPath)

	for _, entry := range entries {

		sfp := path.Join(srcPath, entry.Name())
		dfp := path.Join(dstPath, entry.Name())
		if entry.IsDir() {
			err = this.Fs(sfp, dfp)
			if nil != err {
				return
			}
		} else {
			// perform copy
			err = this.fileCopier.Fs(sfp, dfp)
			if nil != err {
				return
			}
		}

	}
	return
}

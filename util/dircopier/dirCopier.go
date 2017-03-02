package dircopier

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ DirCopier

import (
	"fmt"
	"github.com/opspec-io/opctl/util/filecopier"
	"github.com/opspec-io/opctl/util/vfs"
	"github.com/opspec-io/opctl/util/vfs/os"
	"io/ioutil"
	"path"
)

type DirCopier interface {
	// copies a fs dir from srcPath to dstPath. Creates or overwrites the destination as needed.
	Fs(srcPath string, dstPath string) (err error)
}

func New() DirCopier {
	fs := os.New()
	return dirCopier{
		fs:         fs,
		fileCopier: filecopier.New(),
	}
}

type dirCopier struct {
	fs         vfs.Vfs
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

	// @TODO: remove dependence on real fs here
	entries, err := ioutil.ReadDir(srcPath)

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

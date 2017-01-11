package os

import (
	"github.com/opspec-io/opctl/util/fs"
	"os"
)

func New() fs.Fs {
	return _fs{}
}

type _fs struct{}

func (this _fs) Create(name string) (*os.File, error) {
	return os.Create(name)
}

func (this _fs) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (this _fs) RemoveAll(path string) error {
	return os.RemoveAll(path)
}

func (this _fs) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

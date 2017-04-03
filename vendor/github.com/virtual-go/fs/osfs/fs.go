package osfs

import (
	"github.com/virtual-go/fs"
	"os"
)

func New() fs.FS {
	return _fs{}
}

type _fs struct{}

func (this _fs) Chmod(name string, mode os.FileMode) error {
	return os.Chmod(name, mode)
}

func (this _fs) Create(name string) (*os.File, error) {
	return os.Create(name)
}

func (this _fs) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (this _fs) Open(name string) (*os.File, error) {
	return os.Open(name)
}

func (this _fs) OpenFile(name string, flag int, perm os.FileMode) (*os.File, error) {
	return os.OpenFile(name, flag, perm)
}

func (this _fs) RemoveAll(path string) error {
	return os.RemoveAll(path)
}

func (this _fs) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

package osfs

import (
	"github.com/virtual-go/vfs"
	"os"
)

func New() vfs.Vfs {
	return _vfs{}
}

type _vfs struct{}

func (this _vfs) Chmod(name string, mode os.FileMode) error {
	return os.Chmod(name, mode)
}

func (this _vfs) Create(name string) (*os.File, error) {
	return os.Create(name)
}

func (this _vfs) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (this _vfs) Open(name string) (*os.File, error) {
	return os.Open(name)
}

func (this _vfs) RemoveAll(path string) error {
	return os.RemoveAll(path)
}

func (this _vfs) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

package iioutil

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ IIOUtil

import (
	"io"
	"io/ioutil"
	"os"
)

// virtual filesystem interface
type IIOUtil interface {
	// ReadAll reads from r until an error or EOF and returns the data it read.
	// A successful call returns err == nil, not err == EOF. Because ReadAll is
	// defined to read from src until EOF, it does not treat an EOF from Read
	// as an error to be reported.
	ReadAll(r io.Reader) ([]byte, error)

	// ReadDir reads the directory named by dirname and returns
	// a list of directory entries sorted by filename.
	ReadDir(dirname string) ([]os.FileInfo, error)

	// ReadFile reads the file named by filename and returns the contents.
	// A successful call returns err == nil, not err == EOF. Because ReadFile
	// reads the whole file, it does not treat an EOF from Read as an error
	// to be reported.
	ReadFile(filename string) ([]byte, error)

	// WriteFile writes data to a file named by filename.
	// If the file does not exist, WriteFile creates it with permissions perm;
	// otherwise WriteFile truncates it before writing.
	WriteFile(filename string, data []byte, perm os.FileMode) error
}

func New() IIOUtil {
	return _IIOUtil{}
}

type _IIOUtil struct{}

func (iou _IIOUtil) ReadAll(r io.Reader) ([]byte, error) {
	return ioutil.ReadAll(r)
}

func (iou _IIOUtil) ReadDir(dirname string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(dirname)
}

func (iou _IIOUtil) ReadFile(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}

func (iou _IIOUtil) WriteFile(filename string, data []byte, perm os.FileMode) error {
	return ioutil.WriteFile(filename, data, perm)
}

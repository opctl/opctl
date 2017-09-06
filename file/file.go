// Package file implements usecases surrounding files
package file

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ File

type File interface {
	validator
}

func New() File {
	return _File{
		validator: newValidator(),
	}
}

type _File struct {
	validator
}

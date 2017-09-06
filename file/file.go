// Package file implements usecases surrounding files
package file

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ File

type File interface {
	Validator
}

func New() File {
	return _File{
		Validator: newValidator(),
	}
}

type _File struct {
	Validator
}

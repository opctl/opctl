// Package dir implements usecases surrounding dirs
package dir

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ Dir

type Dir interface {
	validator
}

func New() Dir {
	return _Dir{
		validator: newValidator(),
	}
}

type _Dir struct {
	validator
}

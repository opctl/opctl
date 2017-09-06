// Package dir implements usecases surrounding dirs
package dir

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ Dir

type Dir interface {
	Validator
}

func New() Dir {
	return _Dir{
		Validator: newValidator(),
	}
}

type _Dir struct {
	Validator
}

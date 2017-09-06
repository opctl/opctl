// Package object implements usecases surrounding objects
package object

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ Object

type Object interface {
	Validator
}

func New() Object {
	return _Object{
		Validator: newValidator(),
	}
}

type _Object struct {
	Validator
}

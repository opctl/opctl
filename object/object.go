// Package object implements usecases surrounding objects
package object

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ Object

type Object interface {
	validator
}

func New() Object {
	return _Object{
		validator: newValidator(),
	}
}

type _Object struct {
	validator
}

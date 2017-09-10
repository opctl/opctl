// Package number implements usecases surrounding numbers
package number

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ Number

type Number interface {
	interpreter
}

func New() Number {
	return _Number{
		interpreter: newInterpreter(),
	}
}

type _Number struct {
	interpreter
}

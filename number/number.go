// Package number implements usecases surrounding numbers
package number

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ Number

type Number interface {
	Interpreter
}

func New() Number {
	return _Number{
		Interpreter: newInterpreter(),
	}
}

type _Number struct {
	Interpreter
}

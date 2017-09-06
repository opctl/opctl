// Package number implements usecases surrounding numbers
package number

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ Number

type Number interface {
	Interpreter
	Validator
}

func New() Number {
	return _Number{
		Interpreter: newInterpreter(),
		Validator:   newValidator(),
	}
}

type _Number struct {
	Interpreter
	Validator
}

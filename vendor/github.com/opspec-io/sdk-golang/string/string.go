// Package string implements usecases surrounding strings
package string

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ String

type String interface {
	Interpreter
}

func New() String {
	return _String{
		Interpreter: newInterpreter(),
	}
}

type _String struct {
	Interpreter
}

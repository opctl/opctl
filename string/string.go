// Package string implements usecases surrounding strings
package string

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ String

type String interface {
	interpreter
}

func New() String {
	return _String{
		interpreter: newInterpreter(),
	}
}

type _String struct {
	interpreter
}

// Package expression implements usecases surrounding expressions
package expression

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ Expression

type Expression interface {
	evalToFile
	evalToNumber
	evalToString
}

func New() Expression {
	return _Expression{
		evalToFile:   newEvalToFile(),
		evalToNumber: newEvalToNumber(),
		evalToString: newEvalToString(),
	}
}

type _Expression struct {
	evalToFile
	evalToNumber
	evalToString
}

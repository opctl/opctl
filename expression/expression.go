// Package expression implements usecases surrounding expressions
package expression

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ Expression

type Expression interface {
	evalToDir
	evalToFile
	evalToNumber
	evalToString
}

func New() Expression {
	return _Expression{
		evalToDir:    newEvalToDir(),
		evalToFile:   newEvalToFile(),
		evalToNumber: newEvalToNumber(),
		evalToString: newEvalToString(),
	}
}

type _Expression struct {
	evalToDir
	evalToFile
	evalToNumber
	evalToString
}

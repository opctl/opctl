// Package data implements common use cases involving typed data such as validation & coercion
package data

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ Data

type Data interface {
	coerceToNumber
	coerceToObject
	coerceToString
	validator
}

func New() Data {
	return _Data{
		coerceToNumber: newCoerceToNumber(),
		coerceToObject: newCoerceToObject(),
		coerceToString: newCoerceToString(),
		validator:      newValidator(),
	}
}

type _Data struct {
	coerceToNumber
	coerceToObject
	coerceToString
	validator
}

// Package data implements common use cases involving typed data such as validation & coercion
package data

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ Data

type Data interface {
	coercer
	validator
}

func New() Data {
	return _Data{
		coercer:   newCoercer(),
		validator: newValidator(),
	}
}

type _Data struct {
	coercer
	validator
}

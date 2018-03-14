// Package data implements common use cases involving typed data such as validation & coercion
package data

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ Data

type Data interface {
	validator
}

func New() Data {
	return _Data{
		validator: newValidator(),
	}
}

type _Data struct {
	validator
}

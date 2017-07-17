// Package manifest implements use cases for managing opspec package manifests
package manifest

import (
	"github.com/opspec-io/sdk-golang/model"
)

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ Manifest

type Manifest interface {
	Validator

	// Unmarshal unmarshals the pkg manifest at path
	Unmarshal(
		manifestBytes []byte,
	) (*model.PkgManifest, error)
}

func New() Manifest {
	return _Manifest{
		Validator: newValidator(),
	}
}

type _Manifest struct {
	Validator
}

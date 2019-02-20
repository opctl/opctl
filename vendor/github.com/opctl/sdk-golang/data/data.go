// Package data implements use cases specific to data
package data

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ Data

type Data interface {
	providerFactory

	resolver
}

func New() Data {
	return struct {
		providerFactory
		resolver
	}{
		providerFactory: newProviderFactory(),
		resolver:        newResolver(),
	}
}

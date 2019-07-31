// Package data implements use cases specific to data
package data

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -o ./fake.go --fake-name Fake ./ Data

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

// Package data implements use cases specific to data
package data

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

type Data interface {
	providerFactory

	resolver
}

//counterfeiter:generate -o fakes/data.go . Data
func New() Data {
	return struct {
		providerFactory
		resolver
	}{
		providerFactory: newProviderFactory(),
		resolver:        newResolver(),
	}
}

package auth

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

import (
	"github.com/opctl/opctl/cli/internal/dataresolver"
	"github.com/opctl/opctl/cli/internal/nodeprovider"
)

// Auth exposes the "auth" sub command
//counterfeiter:generate -o fakes/auth.go . Auth
type Auth interface {
	Adder
}

// New returns an initialized "auth" sub command
func New(
	dataResolver dataresolver.DataResolver,
	nodeProvider nodeprovider.NodeProvider,
) Auth {
	return _auth{
		Adder: newAdder(nodeProvider),
	}
}

type _auth struct {
	Adder
}

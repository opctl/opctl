package provider

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

import (
	"context"

	"github.com/opctl/opctl/sdks/go/model"
)

const (
	OpFileName       = "op.yml"
	DotOpspecDirName = ".opspec"
)

// Provider is the interface for something that provides data
//counterfeiter:generate -o fakes/provider.go . Provider
type Provider interface {
	// TryResolve resolves a package from the source.
	//
	// expected errs:
	//  - ErrDataProviderAuthentication on authentication failure
	//  - ErrDataProviderAuthorization on authorization failure
	//  - ErrDataRefResolution on resolution failure
	TryResolve(
		ctx context.Context,
		dataRef string,
	) (model.DataHandle, error)
}

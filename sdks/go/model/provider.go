package model

import (
	"context"
)

const (
	OpFileName       = "op.yml"
	DotOpspecDirName = ".opspec"
)

// DataProvider is the interface for something that provides data
type DataProvider interface {
	Label() string

	// TryResolve resolves a package from the source.
	//
	// expected errs:
	//  - ErrDataProviderAuthentication on authentication failure
	//  - ErrDataProviderAuthorization on authorization failure
	//  - ErrDataRefResolution on resolution failure
	TryResolve(
		ctx context.Context,
		dataRef string,
	) (DataHandle, error)
}

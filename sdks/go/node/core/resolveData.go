package core

import (
	"context"
	"github.com/opctl/opctl/sdks/go/types"
)

// Resolve attempts to resolve data via local filesystem or git
// nil pullCreds will be ignored
//
// expected errs:
//  - ErrDataProviderAuthentication on authentication failure
//  - ErrDataProviderAuthorization on authorization failure
//  - ErrDataRefResolution on resolution failure
func (this _core) ResolveData(
	ctx context.Context,
	dataRef string,
	pullCreds *types.PullCreds,
) (
	types.DataHandle,
	error,
) {
	return this.data.Resolve(
		ctx,
		dataRef,
		this.data.NewFSProvider(),
		this.data.NewGitProvider(this.dataCachePath, pullCreds),
	)
}

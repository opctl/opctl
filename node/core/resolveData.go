package core

import (
	"context"
	"github.com/opctl/sdk-golang/model"
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
	pullCreds *model.PullCreds,
) (
	model.DataHandle,
	error,
) {
	return this.data.Resolve(
		ctx,
		dataRef,
		this.data.NewFSProvider(),
		this.data.NewGitProvider(this.dataCachePath, pullCreds),
	)
}

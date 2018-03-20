package core

import (
	"context"
	"github.com/opspec-io/sdk-golang/model"
)

// Resolve attempts to resolve a pkg via local filesystem or git
// nil pullCreds will be ignored
//
// expected errs:
//  - ErrDataProviderAuthentication on authentication failure
//  - ErrDataProviderAuthorization on authorization failure
//  - ErrDataRefResolution on resolution failure
func (this _core) ResolvePkg(
	ctx context.Context,
	pkgRef string,
	pullCreds *model.PullCreds,
) (
	model.DataHandle,
	error,
) {
	return this.data.Resolve(
		ctx,
		pkgRef,
		this.data.NewFSProvider(),
		this.data.NewGitProvider(this.pkgCachePath, pullCreds),
	)
}

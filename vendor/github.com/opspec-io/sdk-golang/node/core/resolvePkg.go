package core

import (
	"github.com/opspec-io/sdk-golang/model"
)

// Resolve attempts to resolve a pkg via local filesystem or git
// nil pullCreds will be ignored
// returns ErrAuthenticationFailed on authentication failure
// returns ErrPkgNotFound on resolution failure
func (this _core) ResolvePkg(
	pkgRef string,
	pullCreds *model.PullCreds,
) (
	model.PkgHandle,
	error,
) {
	return this.pkg.Resolve(
		pkgRef,
		this.pkg.NewFSProvider(),
		this.pkg.NewGitProvider(this.pkgCachePath, pullCreds),
	)
}

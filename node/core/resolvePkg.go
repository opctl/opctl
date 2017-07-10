package core

import "github.com/opspec-io/sdk-golang/pkg"

// Resolve attempts to resolve a package according to opspec package resolution rules
// nil opts will be ignored
// returns ErrAuthenticationFailed on authentication failure
func (this _core) ResolvePkg(
	pkgRef string,
	opts *pkg.ResolveOpts,
) (
	pkg.Handle,
	error,
) {
	return this.pkg.Resolve(pkgRef, opts)
}

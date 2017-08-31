package pkg

import "github.com/opspec-io/sdk-golang/model"

type Resolver interface {
	// TryResolve attempts to resolve a package from providers in order
  //
  // expected errs:
  //  - ErrPkgPullAuthentication on authentication failure
  //  - ErrPkgPullAuthorization on authorization failure
  //  - ErrPkgNotFound on resolution failure
	Resolve(
		pkgRef string,
		providers ...Provider,
	) (
		model.PkgHandle,
		error,
	)
}

func newResolver() Resolver {
	return _Resolver{}
}

type _Resolver struct{}

func (this _Resolver) Resolve(
	pkgRef string,
	providers ...Provider,
) (
	model.PkgHandle,
	error,
) {
	for _, src := range providers {
		handle, err := src.TryResolve(pkgRef)
		if nil != err {
			return nil, err
		} else if nil != handle {
			return handle, nil
		}
	}

	// if we reached this point resolution failed, return err
	return nil, model.ErrPkgNotFound{}
}

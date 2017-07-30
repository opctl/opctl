package pkg

type Resolver interface {
	// TryResolve attempts to resolve a package from providers in order
	// returns ErrAuthenticationFailed on authentication failure
	// returns ErrPkgNotFound on resolution failure
	Resolve(
		pkgRef string,
		providers ...Provider,
	) (
		Handle,
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
	Handle,
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
	return nil, ErrPkgNotFound{}
}

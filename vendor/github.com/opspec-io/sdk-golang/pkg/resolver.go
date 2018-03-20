package pkg

import (
	"context"
	"github.com/opspec-io/sdk-golang/model"
)

type resolver interface {
	// TryResolve attempts to resolve a package from providers in order
	//
	// expected errs:
	//  - ErrDataProviderAuthentication on authentication failure
	//  - ErrDataProviderAuthorization on authorization failure
	//  - ErrPkgNotFound on resolution failure
	Resolve(
		ctx context.Context,
		pkgRef string,
		providers ...Provider,
	) (
		model.PkgHandle,
		error,
	)
}

func newResolver() resolver {
	return _resolver{}
}

type _resolver struct{}

func (rslv _resolver) Resolve(
	ctx context.Context,
	pkgRef string,
	providers ...Provider,
) (
	model.PkgHandle,
	error,
) {
	for _, src := range providers {
		handle, err := src.TryResolve(ctx, pkgRef)
		if nil != err {
			return nil, err
		} else if nil != handle {
			return handle, nil
		}
	}

	// if we reached this point resolution failed, return err
	return nil, model.ErrPkgNotFound{}
}

package data

import (
	"context"
	"github.com/opctl/sdk-golang/model"
)

type resolver interface {
	// TryResolve attempts to resolve "dataRef" from "providers" in order
	//
	// expected errs:
	//  - ErrDataProviderAuthentication on authentication failure
	//  - ErrDataProviderAuthorization on authorization failure
	//  - ErrDataRefResolution on resolution failure
	Resolve(
		ctx context.Context,
		dataRef string,
		providers ...Provider,
	) (
		model.DataHandle,
		error,
	)
}

func newResolver() resolver {
	return _resolver{}
}

type _resolver struct{}

func (rslv _resolver) Resolve(
	ctx context.Context,
	dataRef string,
	providers ...Provider,
) (
	model.DataHandle,
	error,
) {
	for _, src := range providers {
		handle, err := src.TryResolve(ctx, dataRef)
		if nil != err {
			return nil, err
		} else if nil != handle {
			return handle, nil
		}
	}

	// if we reached this point resolution failed, return err
	return nil, model.ErrDataRefResolution{}
}

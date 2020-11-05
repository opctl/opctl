package data

import (
	"context"

	"github.com/opctl/opctl/sdks/go/model"
)

// Resolve "dataRef" from "providers" in order
//
// expected errs:
//  - ErrDataProviderAuthentication on authentication failure
//  - ErrDataProviderAuthorization on authorization failure
//  - ErrDataRefResolution on resolution failure
func Resolve(
	ctx context.Context,
	dataRef string,
	providers ...model.DataProvider,
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

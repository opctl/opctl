package data

import (
	"context"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/pkg/errors"
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
	var errs []error
	for _, src := range providers {
		handle, err := src.TryResolve(ctx, dataRef)
		if err != nil {
			errs = append(errs, errors.Wrap(err, src.Label()))
		} else if handle != nil {
			return handle, nil
		}
	}

	return nil, ErrDataResolution{
		dataRef: dataRef,
		errs:    errs,
	}
}

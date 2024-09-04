package data

import (
	"context"
	"fmt"

	aggregateError "github.com/opctl/opctl/sdks/go/internal/aggregate_error"
	"github.com/opctl/opctl/sdks/go/model"
)

// Resolve "dataRef" from "providers" in order
//
// expected errs:
//   - ErrDataProviderAuthentication on authentication failure
//   - ErrDataProviderAuthorization on authorization failure
//   - ErrDataRefResolution on resolution failure
func Resolve(
	ctx context.Context,
	dataRef string,
	providers ...model.DataProvider,
) (
	model.DataHandle,
	error,
) {
	var agg aggregateError.ErrAggregate
	for _, src := range providers {
		handle, err := src.TryResolve(ctx, dataRef)
		if err != nil {
			agg.AddError(fmt.Errorf("%s: %w", src.Label(), err))
		} else if handle != nil {
			return handle, nil
		}
	}

	return nil, fmt.Errorf("%w op \"%s\": %w", model.ErrDataUnableToResolve{}, dataRef, agg)
}

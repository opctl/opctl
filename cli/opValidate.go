package main

import (
	"context"

	"github.com/opctl/opctl/cli/internal/dataresolver"
	"github.com/opctl/opctl/sdks/go/opspec"
)

func opValidate(
	ctx context.Context,
	dataResolver dataresolver.DataResolver,
	opPath []string,
	opRef string,
) error {
	opDirHandle, err := dataResolver.Resolve(
		ctx,
		opRef,
		opPath,
		nil,
	)
	if err != nil {
		return err
	}

	return opspec.Validate(
		ctx,
		*opDirHandle.Path(),
	)
}

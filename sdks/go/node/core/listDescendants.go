package core

import (
	"context"
	"fmt"

	"github.com/opctl/opctl/sdks/go/model"
)

func (c core) ListDescendants(
	ctx context.Context,
	req model.ListDescendantsReq,
) (
	[]*model.DirEntry,
	error,
) {
	if req.DataRef == "" {
		return []*model.DirEntry{}, fmt.Errorf(`"" not a valid data ref`)
	}

	dataHandle, err := c.ResolveData(ctx, req.DataRef, req.PullCreds)
	if err != nil {
		return nil, err
	}

	return dataHandle.ListDescendants(ctx)
}

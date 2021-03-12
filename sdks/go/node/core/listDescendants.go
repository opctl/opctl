package core

import (
	"context"

	"github.com/opctl/opctl/sdks/go/model"
)

func (c core) ListDescendants(
	ctx context.Context,
	req model.ListDescendantsReq,
) (
	[]*model.DirEntry,
	error,
) {
	if req.PkgRef == "" {
		return []*model.DirEntry{}, nil
	}

	dataHandle, err := c.ResolveData(ctx, req.PkgRef, req.PullCreds)
	if err != nil {
		return nil, err
	}

	return dataHandle.ListDescendants(ctx)
}

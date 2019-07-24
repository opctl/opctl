package data

import (
	"context"

	"github.com/opctl/opctl/sdks/go/node/api/client"
	"github.com/opctl/opctl/sdks/go/types"
)

func newNodeHandle(
	client client.Client,
	dataRef string,
	pullCreds *types.PullCreds,
) types.DataHandle {
	return nodeHandle{
		client:    client,
		dataRef:   dataRef,
		pullCreds: pullCreds,
	}
}

func (nh nodeHandle) GetContent(
	ctx context.Context,
	contentPath string,
) (
	types.ReadSeekCloser,
	error,
) {
	return nh.client.GetData(
		ctx,
		types.GetDataReq{
			ContentPath: contentPath,
			PkgRef:      nh.dataRef,
			PullCreds:   nh.pullCreds,
		},
	)
}

// nodeHandle allows interacting w/ data sourced from an opspec node
type nodeHandle struct {
	client    client.Client
	dataRef   string
	pullCreds *types.PullCreds
}

func (nh nodeHandle) ListDescendants(
	ctx context.Context,
) (
	[]*types.DirEntry,
	error,
) {
	return nh.client.ListDescendants(
		ctx,
		types.ListDescendantsReq{
			PkgRef:    nh.dataRef,
			PullCreds: nh.pullCreds,
		},
	)
}

func (hn nodeHandle) Path() *string {
	return nil
}

func (nh nodeHandle) Ref() string {
	return nh.dataRef
}

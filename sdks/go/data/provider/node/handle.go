package node

import (
	"context"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/api/client"
)

func newHandle(
	client client.Client,
	dataRef string,
	pullCreds *model.PullCreds,
) model.DataHandle {
	return handle{
		client:    client,
		dataRef:   dataRef,
		pullCreds: pullCreds,
	}
}

func (nh handle) GetContent(
	ctx context.Context,
	contentPath string,
) (
	model.ReadSeekCloser,
	error,
) {
	return nh.client.GetData(
		ctx,
		model.GetDataReq{
			ContentPath: contentPath,
			PkgRef:      nh.dataRef,
			PullCreds:   nh.pullCreds,
		},
	)
}

// handle allows interacting w/ data sourced from an opspec node
type handle struct {
	client    client.Client
	dataRef   string
	pullCreds *model.PullCreds
}

func (nh handle) ListDescendants(
	ctx context.Context,
) (
	[]*model.DirEntry,
	error,
) {
	return nh.client.ListDescendants(
		ctx,
		model.ListDescendantsReq{
			PkgRef:    nh.dataRef,
			PullCreds: nh.pullCreds,
		},
	)
}

func (hn handle) Path() *string {
	return nil
}

func (nh handle) Ref() string {
	return nh.dataRef
}

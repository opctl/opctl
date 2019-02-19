package data

import (
	"context"

	"github.com/opctl/sdk-golang/model"
	"github.com/opctl/sdk-golang/node/api/client"
)

func newNodeHandle(
	client client.Client,
	dataRef string,
	pullCreds *model.PullCreds,
) model.DataHandle {
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

// nodeHandle allows interacting w/ data sourced from an opspec node
type nodeHandle struct {
	client    client.Client
	dataRef   string
	pullCreds *model.PullCreds
}

func (nh nodeHandle) ListDescendants(
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

func (hn nodeHandle) Path() *string {
	return nil
}

func (nh nodeHandle) Ref() string {
	return nh.dataRef
}

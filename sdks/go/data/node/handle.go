package node

import (
	"context"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node"
)

func newHandle(
	opNode node.OpNode,
	dataRef string,
	pullCreds *model.Creds,
) model.DataHandle {
	return handle{
		opNode:    opNode,
		dataRef:   dataRef,
		pullCreds: pullCreds,
	}
}

// handle allows interacting w/ data sourced from an opspec node
type handle struct {
	opNode    node.OpNode
	dataRef   string
	pullCreds *model.Creds
}

func (nh handle) GetContent(
	ctx context.Context,
	contentPath string,
) (
	model.ReadSeekCloser,
	error,
) {
	return nh.opNode.GetData(
		ctx,
		model.GetDataReq{
			ContentPath: contentPath,
			PkgRef:      nh.dataRef,
			PullCreds:   nh.pullCreds,
		},
	)
}

func (nh handle) ListDescendants(
	ctx context.Context,
) (
	[]*model.DirEntry,
	error,
) {
	return nh.opNode.ListDescendants(
		ctx,
		model.ListDescendantsReq{
			PkgRef:    nh.dataRef,
			PullCreds: nh.pullCreds,
		},
	)
}

func (handle) Path() *string {
	return nil
}

func (nh handle) Ref() string {
	return nh.dataRef
}

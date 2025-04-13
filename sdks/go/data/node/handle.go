package node

import (
	"context"
	"path"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node"
)

func newHandle(
	node node.Node,
	dataRef string,
	pullCreds *model.Creds,
) model.DataHandle {
	return handle{
		node:      node,
		dataRef:   dataRef,
		pullCreds: pullCreds,
	}
}

// handle allows interacting w/ data sourced from an opctl node
type handle struct {
	node      node.Node
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
	return nh.node.GetData(
		ctx,
		model.GetDataReq{
			DataRef:   path.Join(nh.dataRef, contentPath),
			PullCreds: nh.pullCreds,
		},
	)
}

func (nh handle) ListDescendants(
	ctx context.Context,
) (
	[]*model.DirEntry,
	error,
) {
	return nh.node.ListDescendants(
		ctx,
		model.ListDescendantsReq{
			DataRef:   nh.dataRef,
			PullCreds: nh.pullCreds,
		},
	)
}

func (nh handle) Ref() string {
	return nh.dataRef
}

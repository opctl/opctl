package lazylocalnode

import (
	"context"

	"github.com/opctl/opctl/cli/internal/nodeprovider"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node"
)

type lazyLocalNode struct {
	nodeProvider nodeprovider.NodeProvider
}

// New returns a new OpNode object that lazily ensures a new local node is started
// before it tries to do stuff.
func New(nodeProvider nodeprovider.NodeProvider) node.OpNode {
	return &lazyLocalNode{
		nodeProvider: nodeProvider,
	}
}

// getAPICore ensures a local opctl node is running before returning
// an api handle to it
func (l lazyLocalNode) getAPICore() (node.OpNode, error) {
	nodeHandle, err := l.nodeProvider.CreateNodeIfNotExists()
	if err != nil {
		return nil, err
	}
	return nodeHandle.APIClient(), nil
}

func (l lazyLocalNode) AddAuth(
	ctx context.Context,
	req model.AddAuthReq,
) error {
	c, err := l.getAPICore()
	if err != nil {
		return err
	}
	return c.AddAuth(ctx, req)
}

func (l lazyLocalNode) GetEventStream(
	ctx context.Context,
	req *model.GetEventStreamReq,
) (
	<-chan model.Event,
	error,
) {
	c, err := l.getAPICore()
	if err != nil {
		return nil, err
	}
	return c.GetEventStream(ctx, req)
}

func (l lazyLocalNode) KillOp(
	ctx context.Context,
	req model.KillOpReq,
) (
	err error,
) {
	c, err := l.getAPICore()
	if err != nil {
		return err
	}
	return c.KillOp(ctx, req)
}

func (l lazyLocalNode) StartOp(
	ctx context.Context,
	req model.StartOpReq,
) (
	rootCallID string,
	err error,
) {
	c, err := l.getAPICore()
	if err != nil {
		return "", err
	}
	return c.StartOp(ctx, req)
}

func (l lazyLocalNode) GetData(
	ctx context.Context,
	req model.GetDataReq,
) (
	model.ReadSeekCloser,
	error,
) {
	c, err := l.getAPICore()
	if err != nil {
		return nil, err
	}
	return c.GetData(ctx, req)
}

func (l lazyLocalNode) ListDescendants(
	ctx context.Context,
	req model.ListDescendantsReq,
) (
	[]*model.DirEntry,
	error,
) {
	c, err := l.getAPICore()
	if err != nil {
		return nil, err
	}
	return c.ListDescendants(ctx, req)
}

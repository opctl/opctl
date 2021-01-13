package lazylocalnode

import (
	"context"

	"github.com/opctl/opctl/cli/internal/nodeprovider"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node"
	"github.com/opctl/opctl/sdks/go/node/api/client"
)

type lazyLocalNode struct {
	nodeProvider nodeprovider.NodeProvider
}

// New returns a new core object that lazily ensures a new local node is started
// before it tries to do stuff.
func New(nodeProvider nodeprovider.NodeProvider) node.OpNode {
	return &lazyLocalNode{
		nodeProvider: nodeProvider,
	}
}

// getAPICore ensures a local opctl node is running before returning
// an api handle to it
func (l lazyLocalNode) getAPICore() (*client.APIClient, error) {
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
	core, err := l.getAPICore()
	if err != nil {
		return err
	}
	return core.AddAuth(ctx, req)
}

func (l lazyLocalNode) GetEventStream(
	ctx context.Context,
	req *model.GetEventStreamReq,
) (
	<-chan model.Event,
	error,
) {
	core, err := l.getAPICore()
	if err != nil {
		return nil, err
	}
	return core.GetEventStream(ctx, req)
}

func (l lazyLocalNode) KillOp(
	ctx context.Context,
	req model.KillOpReq,
) (
	err error,
) {
	core, err := l.getAPICore()
	if err != nil {
		return err
	}
	return core.KillOp(ctx, req)
}

func (l lazyLocalNode) StartOp(
	ctx context.Context,
	req model.StartOpReq,
) (
	rootCallID string,
	err error,
) {
	core, err := l.getAPICore()
	if err != nil {
		return "", err
	}
	return core.StartOp(ctx, req)
}

func (l lazyLocalNode) Liveness(
	ctx context.Context,
) error {
	core, err := l.getAPICore()
	if err != nil {
		return err
	}
	return core.Liveness(ctx)
}

func (l lazyLocalNode) GetData(
	ctx context.Context,
	req model.GetDataReq,
) (
	model.ReadSeekCloser,
	error,
) {
	core, err := l.getAPICore()
	if err != nil {
		return nil, err
	}
	return core.GetData(ctx, req)
}

func (l lazyLocalNode) ListDescendants(
	ctx context.Context,
	req model.ListDescendantsReq,
) (
	[]*model.DirEntry,
	error,
) {
	core, err := l.getAPICore()
	if err != nil {
		return nil, err
	}
	return core.ListDescendants(ctx, req)
}

package core

import (
	"context"
	"github.com/opctl/opctl/sdks/go/types"
)

func (this _core) GetEventStream(
	ctx context.Context,
	req *types.GetEventStreamReq,
) (
	<-chan types.Event,
	<-chan error,
) {

	return this.pubSub.Subscribe(ctx, req.Filter)
}

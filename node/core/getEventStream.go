package core

import (
	"context"
	"github.com/opctl/sdk-golang/model"
)

func (this _core) GetEventStream(
	ctx context.Context,
	req *model.GetEventStreamReq,
) (
	<-chan model.Event,
	<-chan error,
) {

	return this.pubSub.Subscribe(ctx, req.Filter)
}

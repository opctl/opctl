package core

import (
	"context"
	"github.com/opctl/opctl/sdks/go/model"
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

package core

import (
	"context"

	"github.com/opctl/opctl/sdks/go/model"
)

func (this core) GetEventStream(
	ctx context.Context,
	req *model.GetEventStreamReq,
) (
	<-chan model.Event,
	error,
) {

	return this.pubSub.Subscribe(ctx, req.Filter)
}

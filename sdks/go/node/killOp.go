package node

import (
	"context"
	"time"

	"github.com/opctl/opctl/sdks/go/model"
)

func (this core) KillOp(
	ctx context.Context,
	req model.KillOpReq,
) error {
	// killing an op is async
	this.pubSub.Publish(
		model.Event{
			CallKillRequested: &model.CallKillRequested{
				Request: req,
			},
			Timestamp: time.Now().UTC(),
		},
	)
	return nil
}

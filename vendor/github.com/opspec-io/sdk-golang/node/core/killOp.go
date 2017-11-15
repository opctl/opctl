package core

import (
	"github.com/opspec-io/sdk-golang/model"
	"time"
)

func (this _core) KillOp(
	req model.KillOpReq,
) {
	this.pubSub.Publish(
		&model.Event{
			Timestamp: time.Now().UTC(),
			OpKilled: &model.OpKilledEvent{
				RootOpId: req.RootOpId,
			},
		},
	)
}

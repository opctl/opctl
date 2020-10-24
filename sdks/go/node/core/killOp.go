package core

import (
	"time"

	"github.com/opctl/opctl/sdks/go/model"
)

func (this _core) KillOp(
	req model.KillOpReq,
) {
	this.pubSub.Publish(
		model.Event{
			OpKillRequested: &model.OpKillRequested{
				Request: req,
			},
			Timestamp: time.Now().UTC(),
		},
	)
}

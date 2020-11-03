package core

import (
	"time"

	"github.com/opctl/opctl/sdks/go/model"
)

func (this _core) KillOp(
	req model.KillOpReq,
) {
	// killing an op is async
	this.pubSub.Publish(
		model.Event{
			CallKillRequested: &model.CallKillRequested{
				Request: req,
			},
			Timestamp: time.Now().UTC(),
		},
	)
}

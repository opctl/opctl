package core

import (
	"time"

	"github.com/opctl/opctl/sdks/go/model"
)

func (this _core) AddAuth(
	req model.AddAuthReq,
) {
	// killing an op is async
	this.pubSub.Publish(
		model.Event{
			AuthAdded: &model.AuthAdded{
				Auth: model.Auth{
					Creds:     req.Creds,
					Resources: req.Resources,
				},
			},
			Timestamp: time.Now().UTC(),
		},
	)
}

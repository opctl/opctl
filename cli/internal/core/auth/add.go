package auth

import (
	"context"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node"
)

// Adder exposes the "auth add" sub command
type Adder interface {
	Add(
		ctx context.Context,
		resources string,
		username string,
		password string,
	) error
}

// newAdder returns an initialized "auth add" sub command
func newAdder(
	core node.OpNode,
) Adder {
	return _adder{
		core: core,
	}
}

type _adder struct {
	core node.OpNode
}

func (ivkr _adder) Add(
	ctx context.Context,
	resources string,
	username string,
	password string,
) error {
	return ivkr.core.AddAuth(
		ctx,
		model.AddAuthReq{
			Resources: resources,
			Creds: model.Creds{
				Username: username,
				Password: password,
			},
		},
	)
}

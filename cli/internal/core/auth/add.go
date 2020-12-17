package auth

import (
	"context"

	"github.com/opctl/opctl/cli/internal/nodeprovider"
	"github.com/opctl/opctl/sdks/go/model"
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
	nodeProvider nodeprovider.NodeProvider,
) Adder {
	return _adder{
		nodeProvider: nodeProvider,
	}
}

type _adder struct {
	nodeProvider nodeprovider.NodeProvider
}

func (ivkr _adder) Add(
	ctx context.Context,
	resources string,
	username string,
	password string,
) error {
	nodeHandle, err := ivkr.nodeProvider.CreateNodeIfNotExists()
	if nil != err {
		return err
	}

	return nodeHandle.APIClient().AddAuth(
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

package auth

import (
	"context"

	"github.com/opctl/opctl/cli/internal/cliexiter"
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
	)
}

// newAdder returns an initialized "auth add" sub command
func newAdder(
	cliExiter cliexiter.CliExiter,
	nodeProvider nodeprovider.NodeProvider,
) Adder {
	return _adder{
		cliExiter:    cliExiter,
		nodeProvider: nodeProvider,
	}
}

type _adder struct {
	cliExiter    cliexiter.CliExiter
	nodeProvider nodeprovider.NodeProvider
}

func (ivkr _adder) Add(
	ctx context.Context,
	resources string,
	username string,
	password string,
) {
	nodeHandle, createNodeIfNotExistsErr := ivkr.nodeProvider.CreateNodeIfNotExists()
	if nil != createNodeIfNotExistsErr {
		ivkr.cliExiter.Exit(cliexiter.ExitReq{Message: createNodeIfNotExistsErr.Error(), Code: 1})
		return // support fake exiter
	}

	err := nodeHandle.APIClient().AddAuth(
		ctx,
		model.AddAuthReq{
			Resources: resources,
			Creds: model.Creds{
				Username: username,
				Password: password,
			},
		},
	)
	if nil != err {
		ivkr.cliExiter.Exit(cliexiter.ExitReq{Message: err.Error(), Code: 1})
		return // support fake exiter
	}
}

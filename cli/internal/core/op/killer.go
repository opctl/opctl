package op

import (
	"context"

	"github.com/opctl/opctl/cli/internal/cliexiter"
	"github.com/opctl/opctl/cli/internal/nodeprovider"
	"github.com/opctl/opctl/sdks/go/model"
)

// Killer exposes the "op kill" sub command
type Killer interface {
	Kill(
		ctx context.Context,
		opID string,
	)
}

// newKiller returns an initialized "op kill" sub command
func newKiller(
	cliExiter cliexiter.CliExiter,
	nodeProvider nodeprovider.NodeProvider,
) Killer {
	return _killer{
		cliExiter:    cliExiter,
		nodeProvider: nodeProvider,
	}
}

type _killer struct {
	cliExiter    cliexiter.CliExiter
	nodeProvider nodeprovider.NodeProvider
}

func (ivkr _killer) Kill(
	ctx context.Context,
	opID string,
) {
	nodeHandle, createNodeIfNotExistsErr := ivkr.nodeProvider.CreateNodeIfNotExists()
	if nil != createNodeIfNotExistsErr {
		ivkr.cliExiter.Exit(cliexiter.ExitReq{Message: createNodeIfNotExistsErr.Error(), Code: 1})
		return // support fake exiter
	}

	err := nodeHandle.APIClient().KillOp(
		ctx,
		model.KillOpReq{
			OpID:     opID,
			RootOpID: opID,
		},
	)
	if nil != err {
		ivkr.cliExiter.Exit(cliexiter.ExitReq{Message: err.Error(), Code: 1})
		return // support fake exiter
	}
}

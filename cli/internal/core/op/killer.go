package op

import (
	"context"

	"github.com/opctl/opctl/cli/internal/cliexiter"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/api/client"
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
	apiClient client.Client,
	cliExiter cliexiter.CliExiter,
) Killer {
	return _killer{
		apiClient: apiClient,
		cliExiter: cliExiter,
	}
}

type _killer struct {
	apiClient client.Client
	cliExiter cliexiter.CliExiter
}

func (ivkr _killer) Kill(
	ctx context.Context,
	opID string,
) {
	err := ivkr.apiClient.KillOp(
		ctx,
		model.KillOpReq{
			OpID: opID,
		},
	)
	if nil != err {
		ivkr.cliExiter.Exit(cliexiter.ExitReq{Message: err.Error(), Code: 1})
		return // support fake exiter
	}
}

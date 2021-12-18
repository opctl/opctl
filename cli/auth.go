package main

import (
	"context"

	"github.com/opctl/opctl/cli/internal/nodeprovider/local"
	"github.com/opctl/opctl/sdks/go/model"
)

// auth implements "auth" command
func auth(
	ctx context.Context,
	nodeConfig local.NodeConfig,
	addAuthReq model.AddAuthReq,
) error {
	np, err := local.New(nodeConfig)
	if err != nil {
		return err
	}

	node, err := np.CreateNodeIfNotExists(ctx)
	if err != nil {
		return err
	}

	return node.AddAuth(
		ctx,
		addAuthReq,
	)
}

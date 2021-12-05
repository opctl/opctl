package main

import (
	"context"

	"github.com/opctl/opctl/cli/internal/nodeprovider"
	"github.com/opctl/opctl/sdks/go/model"
)

// auth implements "auth" command
func auth(
	ctx context.Context,
	nodeProvider nodeprovider.NodeProvider,
	addAuthReq model.AddAuthReq,
) error {
	node, err := nodeProvider.StartNode(ctx)
	if err != nil {
		return err
	}

	return node.AddAuth(
		ctx,
		addAuthReq,
	)
}

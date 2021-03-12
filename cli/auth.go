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
	node, err := nodeProvider.CreateNodeIfNotExists(ctx)
	if nil != err {
		return err
	}

	return node.AddAuth(
		ctx,
		addAuthReq,
	)
}

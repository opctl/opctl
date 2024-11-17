package oppath

import (
	"context"

	"github.com/opctl/opctl/cli/internal/dataresolver"
	"github.com/opctl/opctl/sdks/go/node"
)

// Get the op path
func Get(
	ctx context.Context,
	currentPath string,
	dataResolver dataresolver.DataResolver,
	node node.Node,
) ([]string, error) {
	opPath := []string{
		currentPath,
	}

	defaultOpsDirRef, err := TryGetDefaultOpsDirRef(
		ctx,
		currentPath,
		dataResolver,
		node,
	)
	if err != nil {
		return nil, err
	}

	if defaultOpsDirRef != nil {
		opPath = append(opPath, *defaultOpsDirRef)
	}

	return opPath, nil
}

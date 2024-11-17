package opspath

import (
	"context"

	"github.com/opctl/opctl/sdks/go/node"
)

func Get(
	ctx context.Context,
	startPath string,
	node node.Node,
) ([]string, error) {

	localOpRef, err := GetLocal()
	if err != nil {
		return nil, err
	}

	opRefs := []string{
		localOpRef,
	}

	cfg, err := getConfig(ctx, startPath, node)
	if err != nil {
		return nil, err
	}

	if cfg.Ops.Dir != nil {
		opRefs = append(opRefs, *cfg.Ops.Dir)
	}

	return opRefs, err

}

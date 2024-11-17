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

	for _, op := range cfg.Ops {
		opRefs = append(opRefs, op.Ref)
	}

	return opRefs, err

}

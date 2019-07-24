package data

import (
	"context"
	"github.com/opctl/opctl/sdks/go/node/api/client"
	"github.com/opctl/opctl/sdks/go/types"
	"net/url"
)

// NewNodeProvider returns a pkg provider which sources pkgs from a node
func (pf _providerFactory) NewNodeProvider(
	apiBaseURL url.URL,
	pullCreds *types.PullCreds,
) Provider {
	return nodeProvider{
		nodeClient: client.New(apiBaseURL, nil),
		puller:     newPuller(),
		pullCreds:  pullCreds,
	}
}

type nodeProvider struct {
	nodeClient client.Client
	puller     puller
	pullCreds  *types.PullCreds
}

func (np nodeProvider) TryResolve(
	ctx context.Context,
	dataRef string,
) (types.DataHandle, error) {

	// ensure resolvable by listing contents w/out err
	if _, err := np.nodeClient.ListDescendants(
		ctx,
		types.ListDescendantsReq{
			PkgRef:    dataRef,
			PullCreds: np.pullCreds,
		},
	); nil != err {
		return nil, err
	}

	return newNodeHandle(np.nodeClient, dataRef, np.pullCreds), nil
}

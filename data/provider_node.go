package data

import (
	"context"
	"github.com/opctl/sdk-golang/model"
	"github.com/opctl/sdk-golang/node/api/client"
	"net/url"
)

// NewNodeProvider returns a pkg provider which sources pkgs from a node
func (pf _providerFactory) NewNodeProvider(
	apiBaseURL url.URL,
	pullCreds *model.PullCreds,
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
	pullCreds  *model.PullCreds
}

func (np nodeProvider) TryResolve(
	ctx context.Context,
	dataRef string,
) (model.DataHandle, error) {

	// ensure resolvable by listing contents w/out err
	if _, err := np.nodeClient.ListDescendants(
		ctx,
		model.ListDescendantsReq{
			PkgRef:    dataRef,
			PullCreds: np.pullCreds,
		},
	); nil != err {
		return nil, err
	}

	return newNodeHandle(np.nodeClient, dataRef, np.pullCreds), nil
}

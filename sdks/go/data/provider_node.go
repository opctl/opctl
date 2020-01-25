package data

import (
	"context"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/api/client"
)

// NewNodeProvider returns a pkg provider which sources pkgs from a node
func (pf _providerFactory) NewNodeProvider(
	apiClient client.Client,
	pullCreds *model.PullCreds,
) Provider {
	return nodeProvider{
		apiClient: apiClient,
		puller:    newPuller(),
		pullCreds: pullCreds,
	}
}

type nodeProvider struct {
	apiClient client.Client
	puller    puller
	pullCreds *model.PullCreds
}

func (np nodeProvider) TryResolve(
	ctx context.Context,
	dataRef string,
) (model.DataHandle, error) {

	// ensure resolvable by listing contents w/out err
	if _, err := np.apiClient.ListDescendants(
		ctx,
		model.ListDescendantsReq{
			PkgRef:    dataRef,
			PullCreds: np.pullCreds,
		},
	); nil != err {
		return nil, err
	}

	return newNodeHandle(np.apiClient, dataRef, np.pullCreds), nil
}

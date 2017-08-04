package pkg

import (
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/node/api/client"
	"net/url"
)

// NewNodeProvider returns a pkg provider which sources pkgs from a node
func (pf _ProviderFactory) NewNodeProvider(
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
	pkgRef string,
) (model.PkgHandle, error) {

	// @TODO: handle not found rather than blindly returning handle
	return newNodeHandle(np.nodeClient, pkgRef, np.pullCreds), nil
}

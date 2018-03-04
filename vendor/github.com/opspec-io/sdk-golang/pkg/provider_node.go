package pkg

import (
	"context"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/node/api/client"
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
	pkgRef string,
) (model.PkgHandle, error) {

	// ensure resolvable by listing contents w/out err
	if _, err := np.nodeClient.ListPkgContents(
		ctx,
		model.ListPkgContentsReq{
			PkgRef:    pkgRef,
			PullCreds: np.pullCreds,
		},
	); nil != err {
		return nil, err
	}

	return newNodeHandle(np.nodeClient, pkgRef, np.pullCreds), nil
}

package node

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

import (
	"context"

	"github.com/opctl/opctl/sdks/go/data/provider"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/api/client"
)

// New returns a data provider which sources pkgs from a node
func New(
	apiClient client.Client,
	pullCreds *model.PullCreds,
) provider.Provider {
	return _node{
		apiClient: apiClient,
		pullCreds: pullCreds,
	}
}

type _node struct {
	apiClient client.Client
	pullCreds *model.PullCreds
}

func (np _node) TryResolve(
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

	return newHandle(np.apiClient, dataRef, np.pullCreds), nil
}

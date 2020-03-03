package data

import (
	"github.com/opctl/opctl/sdks/go/data/provider"
	"github.com/opctl/opctl/sdks/go/data/provider/fs"
	"github.com/opctl/opctl/sdks/go/data/provider/git"
	"github.com/opctl/opctl/sdks/go/data/provider/node"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/api/client"
)

type providerFactory interface {
	// NewFSProvider returns a data provider which sources pkgs from the filesystem
	NewFSProvider(
		basePaths ...string,
	) provider.Provider

	// NewGitProvider returns a data provider which sources pkgs from git repos
	NewGitProvider(
		basePath string,
		pullCreds *model.PullCreds,
	) provider.Provider

	// NewNodeProvider returns a data provider which sources pkgs from a node
	NewNodeProvider(
		apiClient client.Client,
		pullCreds *model.PullCreds,
	) provider.Provider
}

func newProviderFactory() providerFactory {
	return _providerFactory{}
}

func (pf _providerFactory) NewFSProvider(
	basePaths ...string,
) provider.Provider {
	return fs.New(
		basePaths...,
	)
}

func (pf _providerFactory) NewGitProvider(
	basePath string,
	pullCreds *model.PullCreds,
) provider.Provider {
	return git.New(
		basePath,
		pullCreds,
	)
}

func (pf _providerFactory) NewNodeProvider(
	apiClient client.Client,
	pullCreds *model.PullCreds,
) provider.Provider {
	return node.New(
		apiClient,
		pullCreds,
	)
}

type _providerFactory struct{}

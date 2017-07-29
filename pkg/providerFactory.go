package pkg

//go:generate counterfeiter -o ./fakeProvider.go --fake-name FakeProvider ./ Provider

import "github.com/opspec-io/sdk-golang/model"

// Provider is the interface for something that provides pkgs
type Provider interface {
	// TryResolve resolves a package from the source.
	// returns ErrPkgNotFound on failure to find package
	TryResolve(
		pkgRef string,
	) (Handle, error)
}

type ProviderFactory interface {
	NewLocalFSProvider(
		basePaths ...string,
	) Provider

	NewGitProvider(
		basePath string,
		pullCreds *model.PullCreds,
	) Provider
}

func newProviderFactory() ProviderFactory {
	return _ProviderFactory{}
}

type _ProviderFactory struct{}

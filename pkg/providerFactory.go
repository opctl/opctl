package pkg

//go:generate counterfeiter -o ./fakeProvider.go --fake-name FakeProvider ./ Provider

import (
	"github.com/opspec-io/sdk-golang/model"
	"net/url"
)

// Provider is the interface for something that provides pkgs
type Provider interface {
	// TryResolve resolves a package from the source.
  //
  // expected errs:
  //  - ErrPkgPullAuthentication on authentication failure
  //  - ErrPkgPullAuthorization on authorization failure
  //  - ErrPkgNotFound on resolution failure
	TryResolve(
		pkgRef string,
	) (model.PkgHandle, error)
}

type ProviderFactory interface {
	// NewFSProvider returns a pkg provider which sources pkgs from the filesystem
	NewFSProvider(
		basePaths ...string,
	) Provider

	// NewGitProvider returns a pkg provider which sources pkgs from git repos
	NewGitProvider(
		basePath string,
		pullCreds *model.PullCreds,
	) Provider

	// NewNodeProvider returns a pkg provider which sources pkgs from a node
	NewNodeProvider(
		apiBaseURL url.URL,
		pullCreds *model.PullCreds,
	) Provider
}

func newProviderFactory() ProviderFactory {
	return _ProviderFactory{}
}

type _ProviderFactory struct{}

package data

//go:generate counterfeiter -o ./fakeProvider.go --fake-name FakeProvider ./ Provider

import (
	"context"
	"github.com/opctl/sdk-golang/model"
	"net/url"
)

// Provider is the interface for something that provides pkgs
type Provider interface {
	// TryResolve resolves a package from the source.
	//
	// expected errs:
	//  - ErrDataProviderAuthentication on authentication failure
	//  - ErrDataProviderAuthorization on authorization failure
	//  - ErrDataRefResolution on resolution failure
	TryResolve(
		ctx context.Context,
		dataRef string,
	) (model.DataHandle, error)
}

type providerFactory interface {
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

func newProviderFactory() providerFactory {
	return _providerFactory{}
}

type _providerFactory struct{}

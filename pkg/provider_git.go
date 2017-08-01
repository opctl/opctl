package pkg

import (
	"github.com/opspec-io/sdk-golang/model"
	"path/filepath"
)

func (pf _ProviderFactory) NewGitProvider(
	basePath string,
	pullCreds *model.PullCreds,
) Provider {
	return gitProvider{
		localFSProvider: pf.NewFSProvider(basePath),
		basePath:        basePath,
		puller:          newPuller(),
		pullCreds:       pullCreds,
	}
}

type gitProvider struct {
	// composed of fsProvider
	localFSProvider Provider
	basePath        string
	puller          puller
	pullCreds       *model.PullCreds
}

func (grp gitProvider) TryResolve(
	pkgRef string,
) (Handle, error) {

	// attempt to resolve from cache
	handle, err := grp.localFSProvider.TryResolve(pkgRef)
	if nil != err {
		return nil, err
	} else if nil != handle {
		return handle, nil
	}

	// attempt pull if cache miss
	err = grp.puller.Pull(grp.basePath, pkgRef, grp.pullCreds)
	if nil != err {
		return nil, err
	}
	return newGitHandle(filepath.Join(grp.basePath, pkgRef), pkgRef), nil
}

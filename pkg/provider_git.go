package pkg

import (
	"github.com/opspec-io/sdk-golang/model"
	"path/filepath"
)

func (pf _providerFactory) NewGitProvider(
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

func (gp gitProvider) TryResolve(
	pkgRef string,
) (model.PkgHandle, error) {

	// attempt to resolve from cache
	handle, err := gp.localFSProvider.TryResolve(pkgRef)
	if nil != err {
		return nil, err
	} else if nil != handle {
		return handle, nil
	}

	// attempt pull if cache miss
	err = gp.puller.Pull(gp.basePath, pkgRef, gp.pullCreds)
	if nil != err {
		return nil, err
	}
	return newGitHandle(filepath.Join(gp.basePath, pkgRef), pkgRef), nil
}

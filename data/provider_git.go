package data

import (
	"context"
	"github.com/opctl/sdk-golang/model"
	"golang.org/x/sync/singleflight"
	"path/filepath"
)

// singleFlightGroup is used to ensure resolves don't race across provider intances
var resolveSingleFlightGroup singleflight.Group

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
	ctx context.Context,
	dataRef string,
) (model.DataHandle, error) {

	// attempt to resolve within singleFlight.Group to ensure concurrent resolves don't race
	handle, err, _ := resolveSingleFlightGroup.Do(
		dataRef,
		func() (interface{}, error) {
			// attempt to resolve from cache
			handle, err := gp.localFSProvider.TryResolve(ctx, dataRef)
			if nil != err {
				return nil, err
			} else if nil != handle {
				return handle, nil
			}

			// attempt pull if cache miss
			err = gp.puller.Pull(ctx, gp.basePath, dataRef, gp.pullCreds)
			if nil != err {
				return nil, err
			}
			return newGitHandle(filepath.Join(gp.basePath, dataRef), dataRef), nil
		},
	)

	if nil != err {
		return nil, err
	}
	return handle.(model.DataHandle), nil
}

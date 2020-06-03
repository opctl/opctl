package local

import (
	"github.com/golang-interfaces/ios"
	"github.com/golang-utils/lockfile"
	"github.com/opctl/opctl/cli/internal/datadir"
	"github.com/opctl/opctl/cli/internal/model"
	"github.com/opctl/opctl/cli/internal/nodeprovider"
)

// New returns an initialized "local" node provider
func New(
	opts model.NodeCreateOpts,
) nodeprovider.NodeProvider {
	dataDir, newDataDirErr := datadir.New(opts.DataDir)
	if nil != newDataDirErr {
		panic(newDataDirErr)
	}

	return nodeProvider{
		dataDir:       dataDir,
		listenAddress: opts.ListenAddress,
		lockfile:      lockfile.New(),
		os:            ios.New(),
	}
}

type nodeProvider struct {
	dataDir       datadir.DataDir
	listenAddress string
	lockfile      lockfile.LockFile
	os            ios.IOS
}

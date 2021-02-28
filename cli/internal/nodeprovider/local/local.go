package local

import (
	"github.com/golang-utils/lockfile"
	"github.com/opctl/opctl/cli/internal/datadir"
	"github.com/opctl/opctl/cli/internal/nodeprovider"
)

// NodeCreateOpts are options for creating a local opctl node
type NodeCreateOpts struct {
	// DataDir sets the path of dir used to store node data
	DataDir string
	// ListenAddress sets the HOST:PORT on which the node will listen
	ListenAddress    string
	ContainerRuntime string
}

// New returns an initialized "local" node provider
func New(
	opts NodeCreateOpts,
) nodeprovider.NodeProvider {
	dataDir, newDataDirErr := datadir.New(opts.DataDir)
	if nil != newDataDirErr {
		panic(newDataDirErr)
	}

	return nodeProvider{
		dataDir:       dataDir,
		listenAddress: opts.ListenAddress,
		lockfile:      lockfile.New(),
	}
}

type nodeProvider struct {
	dataDir       datadir.DataDir
	listenAddress string
	lockfile      lockfile.LockFile
}

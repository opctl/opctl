package local

import (
	"github.com/golang-utils/lockfile"
	"github.com/opctl/opctl/cli/internal/datadir"
	"github.com/opctl/opctl/cli/internal/nodeprovider"
)

// NodeConfig are options for creating a local opctl node
type NodeConfig struct {
	// DataDir sets the path of dir used to store node data
	DataDir string
	// ListenAddress sets the HOST:PORT on which the node will listen
	ListenAddress    string
	ContainerRuntime string
}

// New returns an initialized "local" node provider
func New(
	config NodeConfig,
) nodeprovider.NodeProvider {
	dataDir, newDataDirErr := datadir.New(config.DataDir)
	if newDataDirErr != nil {
		panic(newDataDirErr)
	}

	return nodeProvider{
		containerRuntime: config.ContainerRuntime,
		dataDir:          dataDir,
		listenAddress:    config.ListenAddress,
		lockfile:         lockfile.New(),
	}
}

type nodeProvider struct {
	containerRuntime string
	dataDir          datadir.DataDir
	listenAddress    string
	lockfile         lockfile.LockFile
}

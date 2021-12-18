package local

import (
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
) (nodeprovider.NodeProvider, error) {
	dataDir, err := datadir.New(config.DataDir)
	if err != nil {
		return nil, err
	}

	return nodeProvider{
		config:  config,
		dataDir: dataDir,
	}, nil
}

type nodeProvider struct {
	config  NodeConfig
	dataDir datadir.DataDir
}

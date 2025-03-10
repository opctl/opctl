package local

import (
	"github.com/opctl/opctl/cli/internal/nodeprovider"
)

// NodeConfig are options for creating a local opctl node
type NodeConfig struct {
	// APIListenAddress sets the IP:PORT on which the API server will listen
	APIListenAddress string
	ContainerRuntime string
	// DataDir sets the path of dir used to store node data
	DataDir string
}

// New returns an initialized "local" node provider
func New(
	config NodeConfig,
) (nodeprovider.NodeProvider, error) {
	return nodeProvider{
		config: config,
	}, nil
}

type nodeProvider struct {
	config NodeConfig
}

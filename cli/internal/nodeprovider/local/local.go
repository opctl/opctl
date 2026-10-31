package local

import (
	"github.com/opctl/opctl/cli/internal/nodeprovider"
	"time"
)

// NodeConfig are options for creating a local opctl node
type NodeConfig struct {
	// APIListenAddress sets the IP:PORT on which the API server will listen
	APIListenAddress string
	ContainerRuntime string
	// DataDir sets the path of dir used to store node data
	DataDir string
	// DNSListenAddress sets the IP:PORT on which the DNS server will listen
	DNSListenAddress string
	// ScratchTTL is how long a finished op's scratch is kept when its caller didn't ask for it to be
	// cleaned (the CLI's run always does); 0 keeps it
	ScratchTTL time.Duration
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

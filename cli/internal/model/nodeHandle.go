package model

import (
	"github.com/opctl/opctl/sdks/go/node/api/client"
)

//counterfeiter:generate -o fakes/nodeHandle.go . NodeHandle
type NodeHandle interface {
	// APIClient returns an API client for this node
	APIClient() client.Client
}

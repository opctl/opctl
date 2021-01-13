package nodeprovider

import (
	"github.com/opctl/opctl/sdks/go/node/api/client"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

//counterfeiter:generate -o fakes/nodeHandle.go . NodeHandle

type NodeHandle interface {
	// APIClient returns an API client for this node
	APIClient() client.APIClient
}

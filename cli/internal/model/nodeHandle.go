package model

import (
	"github.com/opctl/opctl/sdks/go/node/api/client"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -o ./fakeNodeHandle.go --fake-name FakeNodeHandle ./ NodeHandle

type NodeHandle interface {
	// APIClient returns an API client for this node
	APIClient() client.Client
}

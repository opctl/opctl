package auth

import (
	"testing"

	dataResolverFakes "github.com/opctl/opctl/cli/internal/dataresolver/fakes"
	nodeFakes "github.com/opctl/opctl/sdks/go/node/fakes"
)

func TestNewAuth(t *testing.T) {
	New(
		new(dataResolverFakes.FakeDataResolver),
		new(nodeFakes.FakeOpNode),
	)
}

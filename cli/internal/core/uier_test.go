package core

import (
	"testing"

	dataResolverFakes "github.com/opctl/opctl/cli/internal/dataresolver/fakes"
	nodeFakes "github.com/opctl/opctl/sdks/go/node/fakes"
)

func TestNewUIer(t *testing.T) {
	newUIer(
		new(dataResolverFakes.FakeDataResolver),
		new(nodeFakes.FakeOpNode),
	)
}

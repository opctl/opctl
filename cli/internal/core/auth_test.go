package core

import (
	"testing"

	dataResolverFakes "github.com/opctl/opctl/cli/internal/dataresolver/fakes"
	nodeFakes "github.com/opctl/opctl/sdks/go/node/fakes"
)

func TestNewAuther(t *testing.T) {
	if newAuther(
		new(dataResolverFakes.FakeDataResolver),
		new(nodeFakes.FakeOpNode),
	).Auth() == nil {
		t.Error("Auther should provide Auth")
	}
}

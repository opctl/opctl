package core

import (
	"testing"

	dataResolverFakes "github.com/opctl/opctl/cli/internal/dataresolver/fakes"
	nodeFakes "github.com/opctl/opctl/sdks/go/node/fakes"
)

func TestNewOper(t *testing.T) {
	if newOper(
		new(dataResolverFakes.FakeDataResolver),
		new(nodeFakes.FakeOpNode),
	).Op() == nil {
		t.Error("Oper should provide Op")
	}
}

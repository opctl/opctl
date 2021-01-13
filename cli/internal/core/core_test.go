package core

import (
	"os"
	"testing"

	cliOutputFakes "github.com/opctl/opctl/cli/internal/clioutput/fakes"
	"github.com/opctl/opctl/cli/internal/nodeprovider/local"
)

func TestNewCore(t *testing.T) {
	core := New(
		new(cliOutputFakes.FakeCliOutput),
		local.NodeCreateOpts{
			DataDir: os.TempDir(),
		},
	)
	if core.Auth() == nil {
		t.Error("core should provide Auth")
	}
	if core.Op() == nil {
		t.Error("core should provide Op")
	}
}

package core

import (
	"os"
	"testing"

	cliOutputFakes "github.com/opctl/opctl/cli/internal/clioutput/fakes"
	"github.com/opctl/opctl/cli/internal/nodeprovider/local"
)

func TestNewCore(t *testing.T) {
	// arrange
	objectUnderTest := New(
		new(cliOutputFakes.FakeCliOutput),
		local.NodeCreateOpts{
			DataDir: os.TempDir(),
		},
	)

	// assert
	if objectUnderTest.Auth() == nil {
		t.Error("core should provide Auth")
	}
	if objectUnderTest.Op() == nil {
		t.Error("core should provide Op")
	}
}

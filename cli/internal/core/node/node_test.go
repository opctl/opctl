package node

import (
	"testing"

	"github.com/opctl/opctl/cli/internal/nodeprovider/fakes"
)

func TestNodeNew(t *testing.T) {
	New(new(fakes.FakeNodeProvider))
}

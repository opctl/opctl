package local

import (
	"os"
	"testing"
)

func TestNewLocalNode(t *testing.T) {
	New(NodeCreateOpts{
		DataDir: os.TempDir(),
	})
}

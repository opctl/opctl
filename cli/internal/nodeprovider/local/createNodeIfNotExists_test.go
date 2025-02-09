package local

import (
	"context"
	"os"

	. "github.com/onsi/ginkgo"
)

var _ = Context("CreateNodeIfNotExists", func() {
	It("shouldn't panic", func() {
		/* arrange */
		dataDir, err := os.MkdirTemp("", "")
		if err != nil {
			panic(err)
		}

		nodeProvider, err := New(NodeConfig{
			ContainerRuntime: "docker",
			DataDir:          dataDir,
		})
		if err != nil {
			panic(err)
		}

		/* act */
		nodeProvider.CreateNodeIfNotExists(
			context.Background(),
		)
		// note can't assert on result since executing binary is ginkgo not opctl
		// so creation will return error with "flag provided but not defined"
	})
})

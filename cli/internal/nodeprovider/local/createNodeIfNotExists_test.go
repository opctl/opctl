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
		nodeProvider := New(NodeCreateOpts{
			DataDir: dataDir,
		})

		/* act */
		nodeProvider.CreateNodeIfNotExists(
			context.Background(),
		)
	})
})

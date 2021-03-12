package local

import (
	"context"
	. "github.com/onsi/ginkgo"
	"os"
)

var _ = Context("CreateNodeIfNotExists", func() {
	It("shouldn't panic", func() {
		/* arrange */
		nodeProvider := New(NodeCreateOpts{
			DataDir: os.TempDir(),
		})

		/* act */
		nodeProvider.CreateNodeIfNotExists(
			context.Background(),
		)
	})
})

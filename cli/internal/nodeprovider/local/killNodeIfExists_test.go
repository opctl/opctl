package local

import (
	. "github.com/onsi/ginkgo"
	"os"
)

var _ = Context("KillNodeIfExists", func() {
	It("shouldn't panic", func() {
		/* arrange */
		nodeProvider := New(NodeCreateOpts{
			DataDir: os.TempDir(),
		})

		/* act */
		nodeProvider.KillNodeIfExists("")
	})
})

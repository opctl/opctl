package local

import (
	"os"
	. "github.com/onsi/ginkgo"
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

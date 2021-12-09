package local

import (
	"os"

	. "github.com/onsi/ginkgo"
)

var _ = Context("KillNodeIfExists", func() {
	It("shouldn't panic", func() {
		/* arrange */
		dataDir, err := os.MkdirTemp("", "")
		if err != nil {
			panic(err)
		}
		nodeProvider := New(NodeConfig{
			DataDir: dataDir,
		})

		/* act */
		nodeProvider.KillNodeIfExists("")
	})
})

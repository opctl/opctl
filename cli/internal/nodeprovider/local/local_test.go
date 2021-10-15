package local

import (
	"os"

	. "github.com/onsi/ginkgo"
)

var _ = Context("New", func() {
	It("shouldn't panic", func() {
		/* arrange */
		dataDir, err := os.MkdirTemp("", "")
		if err != nil {
			panic(err)
		}

		/* act */
		New(NodeCreateOpts{
			DataDir: dataDir,
		})
	})
})

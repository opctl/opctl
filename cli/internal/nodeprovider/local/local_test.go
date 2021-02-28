package local

import (
	. "github.com/onsi/ginkgo"
	"os"
)

var _ = Context("New", func() {
	It("shouldn't panic", func() {
		/* arrange/act/assert */
		New(NodeCreateOpts{
			DataDir: os.TempDir(),
		})
	})

})

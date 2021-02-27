package local

import (
	"os"
	. "github.com/onsi/ginkgo"
)

var _ = Context("New", func() {
    It("shouldn't panic", func() {
      /* arrange/act/assert */
      New(NodeCreateOpts{
        DataDir: os.TempDir(),
      })
    })
  
})

package local

import (
	"io/ioutil"

	. "github.com/onsi/ginkgo"
)

var _ = Context("KillNodeIfExists", func() {
	It("shouldn't panic", func() {
		/* arrange */
		dataDir, err := ioutil.TempDir("", "")
		if err != nil {
			panic(err)
		}
		nodeProvider := New(NodeCreateOpts{
			DataDir: dataDir,
		})

		/* act */
		nodeProvider.KillNodeIfExists("")
	})
})

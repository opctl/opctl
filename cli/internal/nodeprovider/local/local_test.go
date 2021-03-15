package local

import (
	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Context("New", func() {
	It("shouldn't panic", func() {
		/* arrange */
		dataDir, err := ioutil.TempDir("", "")
		if err != nil {
			panic(err)
		}

		/* act */
		New(NodeCreateOpts{
			DataDir: dataDir,
		})
	})
})

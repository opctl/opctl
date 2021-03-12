package local

import (
	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Context("New", func() {
	It("shouldn't panic", func() {
		/* arrange/act/assert */
		dataDir, err := ioutil.TempDir("", "")
		Expect(err).To(BeNil())
		New(NodeCreateOpts{
			DataDir: dataDir,
		})
	})
})

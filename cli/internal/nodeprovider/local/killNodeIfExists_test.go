package local

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Context("KillNodeIfExists", func() {
	It("shouldn't panic", func() {
		/* arrange */
		dataDir, err := os.MkdirTemp("", "")
		if err != nil {
			panic(err)
		}
		nodeProvider, actualErr := New(NodeConfig{
			DataDir: dataDir,
		})

		/* act */
		nodeProvider.KillNodeIfExists()

		/* assert */
		Expect(actualErr).To(BeNil())
	})
})

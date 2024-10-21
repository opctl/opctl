package local

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Context("New", func() {
	It("shouldn't panic", func() {
		/* arrange */
		dataDir, err := os.MkdirTemp("", "")
		if err != nil {
			panic(err)
		}

		/* act */
		actualNode, actualErr := New(NodeConfig{
			DataDir: dataDir,
		})

		/* assert */
		Expect(actualErr).To(BeNil())
		Expect(actualNode).NotTo(BeNil())
	})
})

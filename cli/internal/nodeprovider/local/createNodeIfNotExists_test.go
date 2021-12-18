package local

import (
	"context"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Context("CreateNodeIfNotExists", func() {
	It("shouldn't panic", func() {
		/* arrange */
		dataDir, err := os.MkdirTemp("", "")
		if err != nil {
			panic(err)
		}

		nodeProvider, err := New(NodeConfig{
			ContainerRuntime: "docker",
			DataDir:          dataDir,
			ListenAddress:    "localhost:42224",
		})
		if err != nil {
			panic(err)
		}

		/* act */
		actualNode, actualErr := nodeProvider.CreateNodeIfNotExists(
			context.Background(),
		)

		/* assert */
		Expect(actualErr).To(BeNil())
		Expect(actualNode).NotTo(BeNil())
	})
})

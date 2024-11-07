package node

import (
	"context"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/opctl/opctl/sdks/go/node/containerruntime/fakes"
)

var _ = Context("core", func() {
	Context("New", func() {
		It("should return Core", func() {
			/* arrange */
			dataDir, err := os.MkdirTemp("", "")
			if err != nil {
				panic(err)
			}

			/* act/assert */
			actual, actualErr := New(
				context.Background(),
				new(FakeContainerRuntime),
				dataDir,
			)

			Expect(actualErr).To(BeNil())
			Expect(actual).To(Not(BeNil()))
		})
	})
})

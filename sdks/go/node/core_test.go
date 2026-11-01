package node

import (
	"context"
	"os"
	"path/filepath"

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

		It("should sweep leftover --no-cache dirs", func() {
			/* arrange */
			dataDir := GinkgoT().TempDir()
			leftoverPath := filepath.Join(dataDir, "nocache", "run", "npm")
			Expect(os.MkdirAll(leftoverPath, 0700)).To(Succeed())

			/* act */
			_, actualErr := New(
				context.Background(),
				new(FakeContainerRuntime),
				dataDir,
			)

			/* assert */
			Expect(actualErr).To(BeNil())
			Expect(filepath.Join(dataDir, "nocache")).NotTo(BeAnExistingFile())
		})
	})
})

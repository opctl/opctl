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
			interruptedSweepPath := filepath.Join(dataDir, "nocache-sweep-interrupted")
			Expect(os.MkdirAll(filepath.Join(interruptedSweepPath, "nocache", "old-run"), 0700)).To(Succeed())
			cachePath := filepath.Join(dataDir, "caches", "npm")
			Expect(os.MkdirAll(cachePath, 0700)).To(Succeed())

			/* act */
			_, actualErr := New(
				context.Background(),
				new(FakeContainerRuntime),
				dataDir,
			)

			/* assert */
			Expect(actualErr).To(BeNil())
			Expect(filepath.Join(dataDir, "nocache")).NotTo(BeAnExistingFile())
			freshPath := filepath.Join(dataDir, "nocache", "fresh-run", "npm")
			Expect(os.MkdirAll(freshPath, 0700)).To(Succeed())
			Eventually(func() []string {
				paths, err := filepath.Glob(filepath.Join(dataDir, "nocache-sweep-*"))
				Expect(err).To(BeNil())
				return paths
			}).Should(BeEmpty())
			Expect(freshPath).To(BeADirectory())
			Expect(cachePath).To(BeADirectory())
		})

		It("should not sweep --no-cache dirs while another node holds the data dir", func() {
			/* arrange */
			dataDir := GinkgoT().TempDir()
			_, err := New(context.Background(), new(FakeContainerRuntime), dataDir)
			Expect(err).To(BeNil())
			livePath := filepath.Join(dataDir, "nocache", "run", "npm")
			Expect(os.MkdirAll(livePath, 0700)).To(Succeed())

			/* act */
			_, actualErr := New(context.Background(), new(FakeContainerRuntime), dataDir)

			/* assert */
			Expect(actualErr).To(HaveOccurred())
			Expect(livePath).To(BeADirectory())
		})
	})
})

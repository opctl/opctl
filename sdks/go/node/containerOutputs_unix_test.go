//go:build unix

package node

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/sys/unix"
)

var _ = Describe("outputs file safe open", func() {
	It("rejects FIFOs without waiting for a writer", func() {
		dir, err := os.MkdirTemp("", "outputs-fifo-")
		Expect(err).NotTo(HaveOccurred())
		defer os.RemoveAll(dir)
		path := filepath.Join(dir, "values.json")
		Expect(unix.Mkfifo(path, 0600)).To(Succeed())
		_, err = readContainerOutputs(path)
		Expect(err).To(MatchError(ContainSubstring("regular file")))
	})
})

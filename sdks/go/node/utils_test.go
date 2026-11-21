package node

import (
	"io"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Container output chunks", func() {
	It("publishes a prompt before newline or EOF", func() {
		reader, writer := io.Pipe()
		defer reader.Close()
		defer writer.Close()
		chunks := make(chan string, 1)
		go readRawChunks(reader, func(chunk []byte) {
			chunks <- string(chunk)
		})

		_, err := io.WriteString(writer, "Name? ")
		Expect(err).NotTo(HaveOccurred())
		Eventually(chunks).Should(Receive(Equal("Name? ")))
	})
})

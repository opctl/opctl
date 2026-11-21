package attach

import (
	"bytes"
	"io"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Terminal newline writer", func() {
	It("adds carriage returns in raw mode and preserves existing CRLF", func() {
		var output bytes.Buffer
		raw := false
		writer := &NewlineWriter{Writer: &output, IsRaw: func() bool { return raw }}

		_, err := io.WriteString(writer, "cooked\n")
		Expect(err).NotTo(HaveOccurred())
		raw = true
		_, err = io.WriteString(writer, "raw\r")
		Expect(err).NotTo(HaveOccurred())
		n, err := io.WriteString(writer, "\nnext\n")
		Expect(err).NotTo(HaveOccurred())
		Expect(n).To(Equal(len("\nnext\n")))
		Expect(output.String()).To(Equal("cooked\nraw\r\nnext\r\n"))
	})
})

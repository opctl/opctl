package node

import (
	"io"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Context("lineEventWriter", func() {
	newWriters := func(published *[]string) (*lineEventWriter, *lineEventWriter) {
		order := new(logOrder)
		return &lineEventWriter{order: order, publish: func(l []byte) { *published = append(*published, "out:"+string(l)) }},
			&lineEventWriter{order: order, publish: func(l []byte) { *published = append(*published, "err:"+string(l)) }}
	}

	It("should publish lines and fragments from both streams in the order they were written", func() {
		/* arrange */
		var published []string
		stdout, stderr := newWriters(&published)

		/* act */
		stdout.Write([]byte("out1\nou"))
		stdout.Write([]byte("t2\nfirst"))
		stderr.Write([]byte("err1\nsecond"))
		stderr.Close()
		stdout.Close()

		/* assert */
		Expect(published).To(Equal([]string{"out:out1\n", "out:out2\n", "out:first", "err:err1\n", "err:second"}))
	})
	It("should keep a long line in one event across many small writes", func() {
		/* arrange */
		var published []string
		stdout, _ := newWriters(&published)
		long := strings.Repeat("x", 5e6)

		/* act */
		for i := 0; i < len(long); i += 1024 {
			stdout.Write([]byte(long[i:min(i+1024, len(long))]))
		}
		stdout.Write([]byte("\n"))

		/* assert */
		Expect(published).To(Equal([]string{"out:" + long + "\n"}))
	})
	It("should reject writes after Close", func() {
		/* arrange */
		var published []string
		stdout, _ := newWriters(&published)
		stdout.Close()

		/* act */
		_, err := stdout.Write([]byte("late\n"))

		/* assert */
		Expect(err).To(Equal(io.ErrClosedPipe))
		Expect(published).To(BeEmpty())
	})
})

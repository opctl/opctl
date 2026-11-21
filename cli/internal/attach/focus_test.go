package attach

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func Test(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CLI container attach")
}

var _ = Describe("Container focus", func() {
	It("releases the terminal even with just one running container", func() {
		var focus focus
		focus.Start("only")
		focus.Cycle()
		Expect(focus.Current()).To(BeEmpty())
		focus.Cycle()
		Expect(focus.Current()).To(Equal("only"))
	})

	It("focuses the latest start, cycles through none, and falls back when the focused container ends", func() {
		var focus focus
		Expect(focus.Current()).To(BeEmpty())
		focus.Start("first")
		focus.Start("second")
		focus.Start("third")
		Expect(focus.Current()).To(Equal("third"))
		focus.Cycle()
		Expect(focus.Current()).To(Equal("second"))
		focus.Cycle()
		Expect(focus.Current()).To(Equal("first"))
		focus.Cycle()
		Expect(focus.Current()).To(BeEmpty())
		focus.Cycle()
		Expect(focus.Current()).To(Equal("third"))
		focus.End("third")
		Expect(focus.Current()).To(Equal("second"))
		focus.End("first")
		Expect(focus.Current()).To(Equal("second"))
		focus.End("second")
		Expect(focus.Current()).To(BeEmpty())
	})
})

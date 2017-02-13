package tcp

import (
	. "github.com/onsi/ginkgo"
	"github.com/opspec-io/opctl/node/core"
)

var _ = Context("tcp", func() {
	Context("New", func() {
		It("should return an instance of Tcp", func() {

			/* arrange */
			var _ = New(new(core.Fake)).(Api)

		})
	})
	Context("Start", func() {
		It("should not panic", func() {

			/* arrange */
			objectUnderTest := New(new(core.Fake))

			/* arrange/act/assert */
			go objectUnderTest.Start()

		})
	})
})

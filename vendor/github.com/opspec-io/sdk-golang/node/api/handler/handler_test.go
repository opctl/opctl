package handler

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/node/core"
)

var _ = Context("handler", func() {
	Context("New", func() {
		It("should return Handler", func() {
			/* arrange/act/assert */
			Expect(New(new(core.Fake))).
				To(Not(BeNil()))
		})
	})
})

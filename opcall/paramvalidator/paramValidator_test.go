package paramvalidator

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("_ParamValidator", func() {
	Context("New()", func() {
		It("should not return nil", func() {
			/* arrange/act/assert */
			Expect(New()).Should(Not(BeNil()))
		})
	})
})

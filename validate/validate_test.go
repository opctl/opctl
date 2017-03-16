package validate

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Param", func() {
	Context("New()", func() {
		It("should not return nil", func() {
			/* arrange/act/assert */
			Expect(New()).Should(Not(BeNil()))
		})
	})
})
